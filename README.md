# thermal

Determination of Heat Gain or Loss and the Surface Temperatures of Insulated Pipe based on C680

```
package thermal // import "github.com/Konstantin8105/thermal"


CONSTANTS

const (
	FlatVerticalSurface Orientation = 1
	FlatHeatFlowDown                = 2
	FlatHeatFlowUp                  = 3
	PipeVertical                    = 1
	PipeHorizontal                  = 2
)

FUNCTIONS

func Cylinder(o io.Writer, Tservice float64, layers []Layer, Tamb float64, es *ExternalSurface, ODpipe float64) (
	Q float64, T []float64, err error)
func Flat(o io.Writer, Tservice float64, layers []Layer, Tamb float64, es *ExternalSurface) (
	Q float64, T []float64, err error)

TYPES

type ExternalSurface struct {
	// Has unexported fields.
}
    ExternalSurface is property of thermal surface

func Emiss(Wind, Emiss float64, orient Orientation) *ExternalSurface

func Surf(Surf float64) *ExternalSurface

type Layer struct {
	Thk float64
	Mat Material
}
    Layer of insulation

type Material interface {
	ConductivityAvg(F1, F2 float64) float64
}
    Material is typical interface with thermal conductivity between 2
    temperatires

func NewMaterialExp(a, b float64) Material
    NewMaterialExp return material with exponents functions

func NewMaterialPolynominal(c ...float64) Material
    NewMaterialPolynominal return material with polynomial thermal conductivity
    property

func NewMaterialType3(
	a1, b1, TL float64,
	a2, b2, TU float64,
	a3, b3 float64,
) Material
    NewMaterialType3 return material type 3

type MaterialExp struct {
	// Has unexported fields.
}
    MaterialExp is material with exponents function thermal conductivity by:

        ln(k) = a + b * T

func (m MaterialExp) ConductivityAvg(F1, F2 float64) float64
    ConductivityAvg return thermal conductivity between 2 temperatires.
    Temperature unit: degree F.

type MaterialPolynomial struct {
	// Has unexported fields.
}
    MaterialPolynomial is material with polynomial thermal conductivity property
    by:

        f[0] + f[1]*T + f[2]*T*T + ...

type MaterialType3 struct {
	// Has unexported fields.
}
    MaterialType3 is material type 3

func (m MaterialType3) ConductivityAvg(F1, F2 float64) float64
    ConductivityAvg return thermal conductivity between 2 temperatires.
    Temperature unit: degree F.

type Orientation int8
```
