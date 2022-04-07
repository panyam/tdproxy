package models

import (
	"fmt"
	// "github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func test_option(price float64, delta float64) *Option {
	return NewOption("TEST", "04-01-2020", fmt.Sprintf("%f", price), true, map[string]interface{}{
		"delta":         delta,
		"open_interest": 10,
		"mark":          1,
	})
}

func create_test_pdist(options ...*Option) *ProbDist {
	return DistFromCalls(options)
}

func prob_range_equals(t *testing.T, probrange ProbRange, x1 float64, x2 float64, lt_prob float64, range_prob float64) bool {
	return math.Abs(x1-probrange.X1) < 0.000001 &&
		math.Abs(x2-probrange.X2) < 0.000001 &&
		math.Abs(lt_prob-probrange.LtProb) < 0.000001 &&
		math.Abs(range_prob-probrange.RangeProb) < 0.000001
}

func TestProbRangeEquals(t *testing.T) {
	/*
		pdist := create_test_pdist(test_option(10, 0.6), test_option(20, 0.3), test_option(30, 0.1), test_option(40, 0.0))
		assert.Equal(t, pdist.ProbLessThan(10), 0.4, "")
		   assert pdist.prob_less_than(10) == 0.4
		   assert pdist.prob_between_prices(10, 40) == 0.6
		   assert prob_range_equals(pdist.distribution[0], 0, 10, 0.4, 0.4)
		   assert prob_range_equals(pdist.distribution[1], 10, 20, 0.7, 0.3)
		   assert prob_range_equals(pdist.distribution[2], 20, 30, 0.9, 0.2)
		   assert prob_range_equals(pdist.distribution[3], 30, 40, 1, 0.1)
		   assert abs(pdist.prob_between_prices(14, 36) - 0.44) < 0.000001
	*/
}

/*
def test_flat_seg_payoff():
    pdist = create_test_pdist((10, 0.6), (20, 0.3), (30, 0.1), (40, 0.0))
    seg = Seg(0, 10, 0, 10)
    assert pdist.calculate_payoff(seg) == {"expected_gain": Fixed(4),
                                           "expected_loss": Fixed(0),
                                           "loss_prob": Fixed(0),
                                           "gain_prob": Fixed(0.4)}
    seg = Seg(0, 20, 0, 10)
    assert pdist.calculate_payoff(seg) == {"expected_gain": Fixed(7),
                                           "expected_loss": Fixed(0),
                                           "loss_prob": Fixed(0),
                                           "gain_prob": Fixed(0.7)}
    # prob between 5 and 25 is 0.6
    seg = Seg(5, 25, 0, 10)
    assert pdist.calculate_payoff(seg) == {"expected_gain": Fixed(6),
                                           "expected_loss": Fixed(0),
                                           "loss_prob": Fixed(0),
                                           "gain_prob": Fixed(0.6)}

def test_flat_seg_payoff():
    pdist = create_test_pdist((10, 0.6), (20, 0.3), (30, 0.1), (40, 0.0))
    seg = Seg(0, 10, 1, -4)
    assert pdist.calculate_payoff(seg) == {"expected_gain": Fixed(.72),
                                           "expected_loss": Fixed(.32),
                                           "gain_prob": Fixed(0.24),
                                           "loss_prob": Fixed(0.16)}

    seg = Seg(0, 20, 1, -4)
    assert pdist.calculate_payoff(seg) == {"expected_gain": Fixed(4.02),
                                           "expected_loss": Fixed(0.32),
                                           "gain_prob": Fixed(0.54),
                                           "loss_prob": Fixed(0.16)}

    seg = Seg(0, 25, 1, -4)
    assert pdist.calculate_payoff(seg) == {"expected_gain": Fixed(5.87),
                                           "expected_loss": Fixed(.32),
                                           "gain_prob": Fixed(0.64),
                                           "loss_prob": Fixed(0.16)}

    seg = Seg(4, 25, 1, -4)
    assert pdist.calculate_payoff(seg) == {"expected_gain": Fixed(5.87),
                                           "expected_loss": Fixed(0),
                                           "gain_prob": Fixed(0.64),
                                           "loss_prob": Fixed(0)}
*/
