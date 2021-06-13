package thermal

import (
	"fmt"
	"os"
)

func ExampleProblem1() {
	Q, T, err := Flat(os.Stdout, 450.0, []Layer{
		Layer{
			Thk: 4.0,
			Mat: NewMaterialExp(-1.62, 0.213e-2),
		},
	}, 10.0, Surf(6.0))
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
	//  SURFACE COEF. USED, BTU/HR.SF.F: 6.00
	//       TOTAL HEAT FLUX, BTU/HR.SF: 36.54
	//
	//  LAYER    INSULATION    CONDUCTIVITY   RESISTANCE TEMPERATURE,F
	//     No  THICKNESS,in  BTU.IN/HR.SF.F  HR.SF.F/BTU  INSIDE OUTSIDE
	//      0          4.00           0.337        11.88  450.00 16.09
}

func ExampleProblem1a() {
	Q, T, err := Flat(os.Stdout, 450.0, []Layer{
		Layer{
			Thk: 4.5,
			Mat: NewMaterialExp(-1.62, 0.213e-2),
		},
	}, 10.0, Surf(6.0))
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
	//  SURFACE COEF. USED, BTU/HR.SF.F: 6.00
	//       TOTAL HEAT FLUX, BTU/HR.SF: 32.51
	//
	//  LAYER    INSULATION    CONDUCTIVITY   RESISTANCE TEMPERATURE,F
	//     No  THICKNESS,in  BTU.IN/HR.SF.F  HR.SF.F/BTU  INSIDE OUTSIDE
	//      0          4.50           0.337        13.37  450.00 15.42
}

func ExampleProblem2() {
	Q, T, err := Cylinder(os.Stdout, 800.0, []Layer{
		Layer{
			Thk: 2.0,
			Mat: NewMaterialPolynominal(0.400, 0.105e-3, 0.286e-6),
		},
	}, 80.0, Surf(1.76), 3.5)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v", err)
		return
	}
	_ = Q
	_ = T

	// Output:
	// HEAT FLOW AND SURFACE TEMPERATURES OF INSULATED EQUPMENT PER ASTM C-680
	//
	//          PIPE OUSIDE INSULLATION: YES
	//         ACTUAL PIPE DIAMETER, IN: 3.5
	//      PIPE SERVICE TEMPERATURE, F: 800.0
	//           AMBIENT TEMPERATURE, F: 80.0
	//  SURFACE COEF. USED, BTU/HR.SF.F: 1.76
	//       TOTAL HEAT FLUX, BTU/HR.SF: 234.80
	//
	//  LAYER    INSULATION    CONDUCTIVITY   RESISTANCE TEMPERATURE,F
	//     No  THICKNESS,in  BTU.IN/HR.SF.F  HR.SF.F/BTU  INSIDE OUTSIDE
	//      0          2.00           0.524         5.45  800.00 147.95
}

func ExampleProblem2a() {
	Q, T, err := Cylinder(os.Stdout, 800.0, []Layer{
		Layer{
			Thk: 2.5,
			Mat: NewMaterialPolynominal(0.400, 0.105e-3, 0.286e-6),
		},
	}, 80.0, Surf(1.76), 3.5)
	if err != nil {
		fmt.Fprintf(os.Stdout, "%v", err)
		return
	}
	_ = Q
	_ = T

	// Output:
	// HEAT FLOW AND SURFACE TEMPERATURES OF INSULATED EQUPMENT PER ASTM C-680
	//
	//          PIPE OUSIDE INSULLATION: YES
	//         ACTUAL PIPE DIAMETER, IN: 3.5
	//      PIPE SERVICE TEMPERATURE, F: 800.0
	//           AMBIENT TEMPERATURE, F: 80.0
	//  SURFACE COEF. USED, BTU/HR.SF.F: 1.76
	//       TOTAL HEAT FLUX, BTU/HR.SF: 205.52
	//
	//  LAYER    INSULATION    CONDUCTIVITY   RESISTANCE TEMPERATURE,F
	//     No  THICKNESS,in  BTU.IN/HR.SF.F  HR.SF.F/BTU  INSIDE OUTSIDE
	//      0          2.50           0.522         7.23  800.00 132.47
}

// func ExampleProblem3() {
// 	Q, T, err := Cylinder(os.Stdout, 800.0, []Layer{
// 		Layer{
// 			Thk: 3.0,
// 			Mat: NewMaterialPolynominal(0.400, 0.105e-3, 0.286e-6),
// 		},
// 	}, 80.0, Emiss(0.0, 0.90, PipeHorizontal), 3.5)
// 	if err != nil {
// 		fmt.Fprintf(os.Stdout, "%v", err)
// 		return
// 	}
// 	_ = Q
// 	_ = T
//
// 	// Output:
// }
