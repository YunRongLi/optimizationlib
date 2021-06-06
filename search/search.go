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

type Method int

const (
	GoldenSection Method = iota
	Fibonacci
)

func (m Method) String() string {
	return [...]string{"GoldenSection", "Fibonacci"}[m]
}

type Searcher struct {
	coe, Eps    float64
	Cost        costFunction
	CurrentMethod Method      
}

func NewSearch(Eps float64, Cost costFunction, m Method) Searcher {
	search := Searcher{}
	search.coe = 1.618
	search.Eps = Eps
	search.Cost = Cost
	search.CurrentMethod = m

	return search
}

func (search Searcher) computeWeights(x, direction []float64, w float64) []float64 {
	wg := make([]float64, len(direction))
	for i, d := range direction {
		wg[i] = x[i] + w * d
	}

	return wg
}

func (search Searcher) phase1(x, direction []float64) Phase1Results {
	var ph1Result Phase1Results
	var delta float64
	delta = 0.1
	ph1Result.g_1 = delta

	ph1Result.fg_1 = search.Cost(search.computeWeights(x, direction, ph1Result.g_1))
	ph1Result.fg_2 = search.Cost(search.computeWeights(x, direction, ph1Result.g_2))

	if ph1Result.fg_1 >= ph1Result.fg_2 {
		return ph1Result
	}

	index := 1
	for {
		if index == 1 {
			ph1Result.g = ph1Result.g_1 + delta * math.Pow(search.coe, float64(index))
			ph1Result.fg = search.Cost(search.computeWeights(x, direction, ph1Result.g))
		} else {
			ph1Result.g_2 = ph1Result.g_1
			ph1Result.g_1 = ph1Result.g
			ph1Result.g = ph1Result.g_1 + delta*math.Pow(search.coe, float64(index))

			ph1Result.fg_2 = ph1Result.fg_1
			ph1Result.fg_1 = ph1Result.fg

			ph1Result.fg = search.Cost(search.computeWeights(x, direction, ph1Result.g))
		}

		if ph1Result.fg_2 > ph1Result.fg_1 && ph1Result.fg_1 < ph1Result.fg {
			return ph1Result
		}

		index++
	}
}

func (search Searcher) gsPhase2(x, direction []float64, ph1Result Phase1Results) float64 {
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
	if interval < search.Eps {
		return (intervalUpper + intervalLower) / 2
	}

	maxIter := 0
	for {
		if math.Pow(0.61893, float64(maxIter)) <= (search.Eps / interval) {
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
				f_beta = search.Cost(search.computeWeights(x, direction, beta))
			} else {
				alpha = intervalLower + rho * interval
				beta  = intervalLower + (1 - rho) * interval
				f_alpha = search.Cost(search.computeWeights(x, direction, alpha))
				f_beta  = search.Cost(search.computeWeights(x, direction, beta)) 
			}
		} else {
			if needChangeBound == needChangeLower {
				alpha = beta
				f_alpha = f_beta
				beta = intervalLower + (1 - rho) * interval
				f_beta = search.Cost(search.computeWeights(x, direction, beta))
			} else if needChangeBound == needChangeUpper {
				beta = alpha
				f_beta = f_alpha
				alpha = intervalLower + rho * interval
				f_alpha = search.Cost(search.computeWeights(x, direction, alpha))
			} else {
				alpha = intervalLower + rho * interval
				beta  = intervalLower + (1 - rho) * interval
				f_alpha = search.Cost(search.computeWeights(x, direction, alpha))
				f_beta  = search.Cost(search.computeWeights(x, direction, beta))
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

func (search Searcher) Search(x, direction []float64) float64 {
	var min_x float64
	ph1Result := Phase1Results{} // no assign specific value, default parameters will be 0
	ph1Result = search.phase1(x, direction)
	switch search.CurrentMethod {
	case GoldenSection:
		min_x = search.gsPhase2(x, direction, ph1Result)
	}

	return min_x
}
