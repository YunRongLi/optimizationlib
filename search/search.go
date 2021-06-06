package search

import (
	"math"
)

type costFunction func([]float64) float64

type Phase1Results struct {
	g_2, g_1, g, fg_2, fg_1, fg float64
}

const (
	noNeedChange = iota
	needChangeUpper
	needChangeLower
	needChangeBoth
)

type FiSearch struct {
	Coe, eps    float64
	cost        costFunction
}

func NewFiSearch(eps float64, cost costFunction) FiSearch {
	search := FiSearch{}
	search.Coe = 1.618
	search.eps = eps
	search.cost = cost

	return search
}

func (search FiSearch) computeWeights(x, direction []float64, w float64) []float64 {
	wg := make([]float64, len(direction))
	for i, d := range direction {
		wg[i] = x[i] + w * d
	}

	return wg
}

func (search FiSearch) phase1(x, direction []float64) Phase1Results {
	var ph1Result Phase1Results
	var delta float64
	delta = 0.1
	ph1Result.g_1 = delta

	ph1Result.fg_1 = search.cost(search.computeWeights(x, direction, ph1Result.g_1))
	ph1Result.fg_2 = search.cost(search.computeWeights(x, direction, ph1Result.g_2))

	if ph1Result.fg_1 >= ph1Result.fg_2 {
		return ph1Result
	}

	index := 1
	for {
		if index == 1 {
			ph1Result.g = ph1Result.g_1 + delta * math.Pow(search.Coe, float64(index))
			ph1Result.fg = search.cost(search.computeWeights(x, direction, ph1Result.g))
		} else {
			ph1Result.g_2 = ph1Result.g_1
			ph1Result.g_1 = ph1Result.g
			ph1Result.g = ph1Result.g_1 + delta*math.Pow(search.Coe, float64(index))

			ph1Result.fg_2 = ph1Result.fg_1
			ph1Result.fg_1 = ph1Result.fg

			ph1Result.fg = search.cost(search.computeWeights(x, direction, ph1Result.g))
		}

		if ph1Result.fg_2 > ph1Result.fg_1 && ph1Result.fg_1 < ph1Result.fg {
			return ph1Result
		}

		index++
	}
}

func (search FiSearch) phase2(x, direction []float64, ph1Result Phase1Results) float64 {
	rho := 0.382
	interval, intervalUpper, intervalLower := 0.0, 0.0, 0.0
	needChangeBound := noNeedChange

	if ph1Result.fg_1 > ph1Result.fg_2 {
		intervalUpper = ph1Result.g_1
	} else {
		intervalUpper = ph1Result.g
	}
	intervalLower = ph1Result.g_2
	
	interval = intervalUpper - intervalLower
	if interval < search.eps {
		return (intervalUpper + intervalLower) / 2
	}

	maxIter := 0
	for {
		if math.Pow(0.61893, float64(maxIter)) <= (search.eps / interval) {
			break;
		}
		maxIter += 1
	}

	alpha, beta, f_alpha, f_beta := 0.0, 0.0, 0.0, 0.0
	for i := 0; i < maxIter; i++ {
		if i == 0 {
			if ph1Result.fg_1 > ph1Result.fg_2 {
				alpha = ph1Result.g_1
				beta = intervalLower + (1 - rho) * interval
				f_alpha = ph1Result.fg_1
				f_beta = search.cost(search.computeWeights(x, direction, beta))
			} else {
				alpha = intervalLower + rho * interval
				beta  = intervalLower + (1 - rho) * interval
				f_alpha = search.cost(search.computeWeights(x, direction, alpha))
				f_beta  = search.cost(search.computeWeights(x, direction, beta)) 
			}
		} else {
			if needChangeBound == needChangeLower {
				alpha = beta
				f_alpha = f_beta
				beta = intervalLower + (1 - rho) * interval
				f_beta = search.cost(search.computeWeights(x, direction, beta))
			} else if needChangeBound == needChangeUpper {
				beta = alpha
				f_beta = f_alpha
				alpha = intervalLower + rho * interval
				f_alpha = search.cost(search.computeWeights(x, direction, alpha))
			} else {
				alpha = intervalLower + rho * interval
				beta  = intervalLower + (1 - rho) * interval
				f_alpha = search.cost(search.computeWeights(x, direction, alpha))
				f_beta  = search.cost(search.computeWeights(x, direction, beta))
			}
		}

		if f_alpha < f_beta {
			intervalUpper = beta
			needChangeBound = needChangeUpper
		} else if f_alpha > f_beta {
			intervalLower = alpha
			needChangeBound = needChangeLower
		} else {
			intervalLower = alpha
			intervalUpper = beta
			needChangeBound = needChangeBoth
		}

		interval = intervalUpper - intervalLower
	}

	return (intervalUpper + intervalLower) / 2
}

func (search FiSearch) Search(x, direction []float64) float64 {
	var min_x float64
	ph1Result := Phase1Results{} // no assign specific value, default parameters will be 0
	ph1Result = search.phase1(x, direction)
	min_x = search.phase2(x, direction, ph1Result)

	return min_x
}
