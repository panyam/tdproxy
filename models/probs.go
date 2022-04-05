package models

type ProbDist struct {
	Distribution []*ProbRange
}

func (pd *ProbDist) IsEmpty() bool {
	return len(pd.Distribution) == 0
}

func (pd *ProbDist) Size() int {
	return len(pd.Distribution)
}

func (pd *ProbDist) LastRange() *ProbRange {
	l := len(pd.Distribution)
	if l == 0 {
		return nil
	} else {
		return pd.Distribution[l-1]
	}
}

type ProbRange struct {
	X1        float64
	X2        float64
	LtProb    float64
	RangeProb float64
}

func DistFromCalls(calls []*Option) *ProbDist {
	out := &ProbDist{}
	for _, opt := range calls {
		if opt.OpenInterest <= 0 || !opt.HasValidDelta() {
			continue
		}
		lt_prob := 1 - opt.Delta
		if out.Size() > 0 && lt_prob < out.LastRange().LtProb {
			// "Delta in calls must be decreasing"
			continue
		}
		last_dist := out.LastRange()
		range_prob := lt_prob
		x1 := 0.0
		if !out.IsEmpty() {
			range_prob = lt_prob - last_dist.LtProb
			x1 = last_dist.X2
		}
		if range_prob <= 0 {
			continue
		}
		next := &ProbRange{X1: x1, X2: opt.StrikePrice, LtProb: lt_prob, RangeProb: range_prob}
		out.Distribution = append(out.Distribution, next)
		if opt.Delta == 0 {
			// Only add one price point at 0 - rest are irrelevant
			break
		}
	}
	return out
}
