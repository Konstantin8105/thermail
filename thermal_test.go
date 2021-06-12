package thermal

import (
	"fmt"
	"math"
	"os"
)

func ExampleProblem1() {
	Q, T, err := Flat(os.Stdout, 450.0, []Layer{
		Layer{
			Thk: 4.0,
			Conductivity: func(F float64) float64 {
				return math.Exp(-1.62 + 0.213e-2*F)
			},
		},
	}, 10.0, 6.0)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v", err)
		return
	}
	_ = Q
	_ = T

	// Output:
	// HEAT FLOW AND SURFACE TEMPERATURES OF INSULATED EQUPMENT PER ASTM C-680
	//
	//  EQUPMENT SERVICE TEMPERATURE, F: 450.0
	//           AMBIENT TEMPERATURE, F: 10.0
	//  SURFACE COEF. USED, BTU/HR.SF.F: 6.0
	//       TOTAL HEAT FLUX, BTU/HR.SF: 36.57
	//
	//  LAYER    INSULATION    CONDUCTIVITY   RESISTANCE TEMPERATURE,F
	//     No  THICKNESS,in  BTU.IN/HR.SF.F  HR.SF.F/BTU  INSIDE OUTSIDE
	//      0          4.00           0.337        11.87  450.00 16.09
}

func ExampleProblem1a() {
	Q, T, err := Flat(os.Stdout, 450.0, []Layer{
		Layer{
			Thk: 4.5,
			Conductivity: func(F float64) float64 {
				return math.Exp(-1.62 + 0.213e-2*F)
			},
		},
	}, 10.0, 6.0)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v", err)
		return
	}
	_ = Q
	_ = T

	// Output:
	// HEAT FLOW AND SURFACE TEMPERATURES OF INSULATED EQUPMENT PER ASTM C-680
	//
	//  EQUPMENT SERVICE TEMPERATURE, F: 450.0
	//           AMBIENT TEMPERATURE, F: 10.0
	//  SURFACE COEF. USED, BTU/HR.SF.F: 6.0
	//       TOTAL HEAT FLUX, BTU/HR.SF: 32.53
	//
	//  LAYER    INSULATION    CONDUCTIVITY   RESISTANCE TEMPERATURE,F
	//     No  THICKNESS,in  BTU.IN/HR.SF.F  HR.SF.F/BTU  INSIDE OUTSIDE
	//      0          4.50           0.337        13.36  450.00 15.42
}
