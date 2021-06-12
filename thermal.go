package thermal

import (
	"fmt"
	"math"
)

func lambda(F float64) float64 {
	return math.Exp(-1.62 + 0.213e-2*F)
}

func Rep1() {

	Tamb := 10.0 // F
	Surf := 6.0  //1.65

	I := 0
	I = I + 1

	Nlayer := 1

	Thk := []float64{4.0}

	T := []float64{450.0, 0}
	K := make([]float64, len(T)-1)

	Thktot := 0.0
	for I := 0; I < Nlayer; I++ {
		Thktot += Thk[I]
	}

	Tdelt := T[0] - Tamb
	for I := 0; I < Nlayer; I++ {
		T[I+1] = T[I] - Thk[I]/Thktot*Tdelt
	}

	var iter int64

T220:
iter ++
	//Ts := T[Nlayer+1]

	var Rs float64
	if Surf > 0 {
		Rs = 1 / Surf
		//Surfc = Surf
	} else {
		//RS = surcof(4, TS, Tamb, Emiss, Wind, Nor, 2.0)
		//Surfc = 1 / RS
	}
	//Kcurve := func() float64{
	Tmiddle := ((T[0] + T[1]) / 2.0)
	K[0] = lambda(Tmiddle)
	//	}()
	Rsum := Rs
	R := make([]float64, len(T))
	for i := 0; i < Nlayer; i++ {
		if K[i] > 0.01 {
			R[i] = Thk[i] / K[i]
			Rsum += R[i]
		}
	}

	Q := (T[0] - Tamb) / Rsum

	Tsum := 0.0
	Tint := T[0]
	for i := 0; i < Nlayer; i++ {
		Tint = T[i] - Q*R[i]
		Tsum += math.Abs(T[i+1] - Tint)
		T[i+1] = Tint
	}

	if Tsum > 1e-5 {
		goto T220
	}

	fmt.Printf("Amount iteration:\t%d\n", iter)
	fmt.Printf("Rs\t%.3f\n", Rs)
	fmt.Printf("R\t%.3f\t%.3f\n", R[0], R[1])
	fmt.Printf("Rsum\t%.3f\n", Rsum)
	fmt.Printf("T\t%.3f\t%.3f\n", T[0], T[1])
	fmt.Printf("Q\t%.3f\n", Q)
}
