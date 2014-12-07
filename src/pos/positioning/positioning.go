package positioning

import (
    ."math"
)

func sgn(x float64) float64 {
	if x >= 0.0 {
		return 1.0
	} else {
		return -1.0
	}
}

func fabs(x float64) float64 {
	if x >= 0.0 {
		return x
	} else {
		return -x
	}
}

func Solve_2d(
	reciever [][]float64,
	pseudolites [][]float64,
	pranges1 float64,
	pranges2 float64) {

	var origin [2]float64
	var len float64
	var tan_theta float64
	var sin_theta float64
	var cos_theta float64
	var d1 float64
	var h1 float64
	var invrotation [2][2]float64
	var pranges [2]float64

	//fmt.Printf("\nPseudolites\n%f %f %f %f\n",
		//pseudolites[0][0],
		//pseudolites[0][1],
		//pseudolites[1][0],
		//pseudolites[1][1])
	//fmt.Printf("\npr1 %f pr2 %f\n", pranges1, pranges2)

	pranges[0] = pranges1
	pranges[1] = pranges2

	origin[0] = pseudolites[0][0]
	origin[1] = pseudolites[0][1]

	pseudolites[0][0] = 0
	pseudolites[0][1] = 0
	pseudolites[1][0] = pseudolites[1][0] - origin[0]
	pseudolites[1][1] = pseudolites[1][1] - origin[1]

	len = Sqrt(Pow(pseudolites[1][0], 2) + Pow(pseudolites[1][1], 2))

	tan_theta = pseudolites[1][1] / pseudolites[1][0]
	cos_theta = sgn(pseudolites[1][0]) / Sqrt(Pow(tan_theta, 2)+1)
	sin_theta = sgn(pseudolites[1][1]) * fabs(tan_theta) / Sqrt(Pow(tan_theta, 2)+1)

	invrotation[0][0] = cos_theta
	invrotation[0][1] = -sin_theta
	invrotation[1][0] = sin_theta
	invrotation[1][1] = cos_theta

	d1 = ((Pow(pranges[0], 2)-Pow(pranges[1], 2))/len + len) / 2

	h1 = Sqrt(Pow(pranges[0], 2) - Pow(d1, 2))

	reciever[0][0] = d1
	reciever[0][1] = h1
	reciever[1][0] = d1
	reciever[1][1] = -h1

	reciever[0][0] = invrotation[0][0]*d1 + invrotation[0][1]*h1
	reciever[0][1] = invrotation[1][0]*d1 + invrotation[1][1]*h1
	reciever[0][0] += origin[0]
	reciever[0][1] += origin[1]

	reciever[1][0] = invrotation[0][0]*d1 + invrotation[0][1]*-h1
	reciever[1][1] = invrotation[1][0]*d1 + invrotation[1][1]*-h1
	reciever[1][0] += origin[0]
	reciever[1][1] += origin[1]
}
