package helper

import (
	"math"
	"sort"

	tk "github.com/eaciit/toolkit"
)

func GetCorrelation(Data1, Data2 tk.M) float64 {
	x := []float64{}
	y := []float64{}
	keys := []string{}

	for k, _ := range Data1 {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, _k := range keys {
		if Data1.Has(_k) && Data2.Has(_k) {
			x = append(x, Data1.GetFloat64(_k))
			y = append(y, Data2.GetFloat64(_k))
		}
	}

	if len(x) == 0 {
		return -1
	}

	_stdx := GetStandardDev(x)
	_stdy := GetStandardDev(y)
	_covxy := GetCovariance(x, y)

	// tk.Printfn("%.2f / ( %.2f * %.2f )", _covxy, _stdx, _stdy)
	// tk.Printfn("Result : %.2f", tk.Div(_covxy, (_stdx*_stdy)))

	return tk.Div(_covxy, (_stdx * _stdy))
}

func GetStandardDev(af []float64) float64 {
	_ret := float64(0)

	_c := float64(len(af))
	_s := float64(0)
	_s2 := float64(0)

	for _, f := range af {
		_s += f
		_s2 += (f * f)
	}

	_ret = tk.Div(((_c * _s2) - (_s * _s)), (_c * _c))
	return math.Sqrt(_ret)
}

func GetCovariance(x, y []float64) float64 {
	mx := float64(0)
	my := float64(0)

	for _, v := range x {
		mx += v
	}

	for _, v := range y {
		my += v
	}

	mx = mx / float64(len(x))
	my = my / float64(len(y))

	uxy := float64(0)
	for i := 0; i < len(x); i++ {
		uxy += (x[i] - mx) * (y[i] - my)
	}

	return tk.Div(uxy, float64(len(x)))
}
