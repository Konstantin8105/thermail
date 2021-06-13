package thermal

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"text/tabwriter"

	"github.com/Konstantin8105/pow"
)

type Material interface {
	// Conductivity(F float64) float64
	ConductivityAvg(F1, F2 float64) float64
}

type MaterialPolynominal struct {
	factors []float64
}

// func (m MaterialPolynominal) Conductivity(F float64) float64 {
// 	K := make([]float64, len(m.factors))
// 	for i := range m.factors {
// 		K[i] = m.factors[i]
// 	}
// 	for i := 1; i < len(m.factors); i++ {
// 		for j := i; j < len(m.factors); j++ {
// 			K[j] *= F
// 		}
// 	}
// 	Ksum := 0.0
// 	for i := range K {
// 		Ksum += K[i]
// 	}
// 	return Ksum
// }

func (m MaterialPolynominal) ConductivityAvg(F2, F1 float64) float64 {
	K := make([]float64, len(m.factors))
	for i := range m.factors {
		K[i] = m.factors[i]
	}
	for i := 0; i < len(m.factors); i++ {
		K[i] *= (pow.En(F2, i+1) - pow.En(F1, i+1)) / float64(i+1)
	}
	Ksum := 0.0
	for i := range K {
		Ksum += K[i]
	}
	Ksum = Ksum / (F2 - F1)
	return Ksum
}

func NewMaterialPolynominal(c ...float64) Material {
	return MaterialPolynominal{factors: c}
}

type MaterialExp struct {
	a, b float64
}

func (m MaterialExp) ConductivityAvg(F1, F2 float64) float64 {
	return 1.0 / (F2 - F1) * (math.Exp(m.a+m.b*F2) - math.Exp(m.a+m.b*F1)) / m.b
}

func NewMaterialExp(a, b float64) Material {
	return MaterialExp{a: a, b: b}
}

type Layer struct {
	Thk float64
	Mat Material
}

// func integral(F1, F2 float64, c func(float64) float64) float64 {
// 	var (
// 		amount = 10000
// 		K      float64
// 		dF     = (F2 - F1) / float64(amount-1)
// 	)
// 	Ki := make([]float64, amount)
// 	for i := 0; i < amount; i++ {
// 		Ki[i] = c(F1 + float64(i)*dF)
// 	}
// 	for i := 0; i < amount-1; i++ {
// 		K += (Ki[i] + Ki[i+1]) / 2.0 * dF
// 	}
// 	K = K / (F2 - F1)
// 	return K
// }

type ExternalSurface struct {
	isSurf bool
	surf   float64
	wind   float64
	emiss  float64
	orient Orientation
}

func Surf(Surf float64) *ExternalSurface {
	es := new(ExternalSurface)
	es.surf = Surf
	es.isSurf = true
	return es
}

func Emiss(Wind, Emiss float64, orient Orientation) *ExternalSurface {
	es := new(ExternalSurface)
	es.wind = Wind
	es.emiss = Emiss
	es.orient = orient
	return es
}

func (es *ExternalSurface) surcof(Dia, Ts, Tamb float64, isCylinder bool) {
	Tair := (Tamb+Ts)/2.0 + 459.69
	ATdelt := math.Abs(Tamb - Ts)
	if ATdelt < 1.0 {
		ATdelt = 1.0
	}

	var Dx float64
	var coef float64

	if isCylinder {
		Dx = Dia * 12.0
		switch es.orient {
		case 1:
			coef = 1.016
		case 2:
			coef = 1.235
		}
		if 24 < Dx {
			Dx = 24.0
		}
	} else {
		Dx = 24.0
		switch es.orient {
		case 1:
			coef = 1.394
		case 2:
			coef = 0.89
		case 3:
			coef = 1.79
		}
	}

	var H, Hsamb, Hramb float64
	Hsamb = coef * math.Pow(Dx, -0.2) * math.Pow(Tair, -0.181) * math.Pow(ATdelt, 0.266) * math.Sqrt(1+1.277*es.wind)
	if Tamb != Ts {
		Hramb = es.emiss * 0.1713e-8 * (pow.E4(Tamb+459.69) - pow.E4(Ts+459.69)) / (Tamb - Ts)
	}
	H = Hsamb + Hramb
	if H < 0.0 {
		H = 1.61
	}

	es.surf = H
}

// NOR = 1 vertical pipe
// NOR = 2 horizontal PIPE

// NOR = 1 horizontal heat FLOW
// NOR = 2 heat flow down
// NOR = 3 heat flow up

type Orientation int8

const (
	FlatVerticalSurface Orientation = 1
	FlatHeatFlowDown                = 2
	FlatHeatFlowUp                  = 3
	PipeVertical                    = 1
	PipeHorizontal                  = 2
)

func Flat(o io.Writer, Tservice float64, layers []Layer, Tamb float64, es *ExternalSurface) (
	Q float64, T []float64, err error) {
	return calc(o, Tservice, layers, Tamb, es, -1.0)
}

func Cylinder(o io.Writer, Tservice float64, layers []Layer, Tamb float64, es *ExternalSurface, ODpipe float64) (
	Q float64, T []float64, err error) {
	return calc(o, Tservice, layers, Tamb, es, ODpipe)
}

func calc(o io.Writer, Tservice float64, layers []Layer, Tamb float64, es *ExternalSurface, ODpipe float64) (
	Q float64, T []float64, err error) {

	isCylinder := 0.0 < ODpipe

	// nil output
	if o == nil {
		var buf bytes.Buffer
		o = &buf
	}
	out := tabwriter.NewWriter(o, 0, 0, 1, ' ', tabwriter.AlignRight)
	defer func() {
		out.Flush()
	}()

	{
		// input data
		fmt.Fprintf(out, "HEAT FLOW AND SURFACE TEMPERATURES OF INSULATED EQUPMENT PER ASTM C-680\n")
		fmt.Fprintf(out, "\n")
		if isCylinder {
			fmt.Fprintf(out, "PIPE OUSIDE INSULLATION:\t YES\n")
			fmt.Fprintf(out, "ACTUAL PIPE DIAMETER, IN:\t %.1f\n", ODpipe)
			fmt.Fprintf(out, "PIPE SERVICE TEMPERATURE, F:\t %.1f\n", Tservice)
		} else {
			fmt.Fprintf(out, "EQUPMENT SERVICE TEMPERATURE, F:\t %.1f\n", Tservice)
		}
		fmt.Fprintf(out, "AMBIENT TEMPERATURE, F:\t %.1f\n", Tamb)
	}

	// calculate diameters per layers
	OD := make([]float64, len(layers))
	ID := make([]float64, len(layers))
	{
		ID[0] = ODpipe
		for i := range layers {
			if 0 < i {
				ID[i] = OD[i-1]
			}
			OD[i] = ID[i] + 2.0*layers[i].Thk
		}
	}

	// temperature initialization
	T = make([]float64, len(layers)+1)
	R := make([]float64, len(layers))
	K := make([]float64, len(layers))
	{
		ThkSum := 0.0
		for _, l := range layers {
			ThkSum += l.Thk
		}
		Tdelta := Tservice - Tamb
		for i := range T {
			if i == 0 {
				T[0] = Tservice
				continue
			}
			T[i] = T[i-1] - layers[i-1].Thk/ThkSum*Tdelta
		}
	}

	var iter, iterMax int64 = 0, 2000
	for ; iter < iterMax; iter++ {
		// symmary
		var Rsum float64
		if !es.isSurf {
			es.surcof(OD[len(layers)-1], T[len(layers)], Tamb, isCylinder)
		}
		Rsum = 1.0 / es.surf
		for i := range layers {
			K[i] = layers[i].Mat.ConductivityAvg(T[i], T[i+1])
			if isCylinder {
				R[i] = OD[len(layers)-1] / 2.0 * math.Log(OD[i]/ID[i]) / K[i]
			} else {
				R[i] = layers[i].Thk / K[i]
			}
			Rsum += R[i]
		}

		// heat flux
		Q = (Tservice - Tamb) / Rsum

		// iteration criteria
		tol := 0.0
		for i := range layers {
			Ts := T[i] - Q*R[i]
			tol += math.Abs(T[i+1] - Ts)
			T[i+1] = Ts // store data
		}
		if math.Abs(tol) < 1e-5 {
			break
		}
	}
	if iterMax <= iter {
		err = fmt.Errorf("not enougnt iterations")
		return
	}

	if isCylinder {
		Q = Q * math.Pi * OD[len(layers)-1] / 12.0
	}

	// TODO orientation view
	{
		// output data
		if !es.isSurf {
			fmt.Fprintf(out, "EMITTANCE:\t %.1f\n", es.emiss)
			fmt.Fprintf(out, "WIND SPEED, MPH:\t %.1f\n", es.wind)
		}
		fmt.Fprintf(out, "SURFACE COEF. USED, BTU/HR.SF.F:\t %.2f\n", es.surf)
		fmt.Fprintf(out, "TOTAL HEAT FLUX, BTU/HR.SF:\t %.2f\n", Q)
		fmt.Fprintf(out, "\n")
		fmt.Fprintf(out, "LAYER \tINSULATION \tCONDUCTIVITY \tRESISTANCE \tTEMPERATURE,F\n")
		fmt.Fprintf(out, "No \tTHICKNESS,in \tBTU.IN/HR.SF.F \tHR.SF.F/BTU \tINSIDE \tOUTSIDE\n")
		for i, l := range layers {
			fmt.Fprintf(out, "%d \t%.2f \t%.3f \t%.2f \t%.2f \t%.2f\n",
				i, l.Thk, K[i], R[i], T[i], T[i+1])
		}
	}
	return
}
