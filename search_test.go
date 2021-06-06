package optimizationlib

import (
	"testing"
	"math"
	"fmt"
)

func TestSearch(t *testing.T) {
	x := []float64{0.1}
	direction := []float64{1}
	var eps float64 = 0.03
	var min_y float64

	test_cost_func := func(x []float64) (float64) {
		var y float64
		y = math.Pow(x[0], 4) - 14 * math.Pow(x[0], 3) + 
			60 * math.Pow(x[0], 2) - 70 * x[0]
		fmt.Println(x[0], y)
		return y
	}

	fiSearch := newFiSearch(eps, test_cost_func)
	x[0] = fiSearch.search(x, direction)
	min_y = test_cost_func(x)
	fmt.Println(x[0], min_y)
}