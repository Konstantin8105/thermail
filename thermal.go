package thermal

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"text/tabwriter"
)

type Layer struct {
	Thk          float64
	Conductivity func(F float64) float64
}

func (l Layer) Kavg(F1, F2 float64) float64 {
	var (
		amount = 10
		K      float64
		dF     = (F2 - F1) / float64(amount-1)
	)
	Ki := make([]float64, amount)
	for i := 0; i < amount; i++ {
		Ki[i] = l.Conductivity(F1 + float64(i)*dF)
	}
	for i := 0; i < amount-1; i++ {
		K += (Ki[i] + Ki[i+1]) / 2.0 * dF
	}
	K = K / (F2 - F1)
	return K
}

func Flat(o io.Writer, Tservice float64, layers []Layer, Tamb float64, Surf float64) (
	Q float64, T []float64, err error) {

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
		fmt.Fprintf(out, "EQUPMENT SERVICE TEMPERATURE, F:\t %.1f\n", Tservice)
		fmt.Fprintf(out, "AMBIENT TEMPERATURE, F:\t %.1f\n", Tamb)
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
		if 0 < Surf {
			Rsum = 1.0 / Surf
		} else {
			// Rsum = 1.0 / surcof(4, TS, Tamb, Emiss, Wind, Nor, 2.0)
		}
		for i := range layers {
			K[i] = layers[i].Kavg(T[i], T[i+1])
			R[i] = layers[i].Thk / K[i]
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

	{
		// output data
		if 0 < Surf {
			fmt.Fprintf(out, "SURFACE COEF. USED, BTU/HR.SF.F:\t %.1f\n", Surf)
		}
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
