package api

import (
	"math"
	"time"
)

// WithConstantDistribution distributes the rate constantly across 100ms intervals
func WithConstantDistribution(iterationDuration time.Duration, rateFn RateFunction) (time.Duration, RateFunction) {
	distributedIterationDuration := 100 * time.Millisecond

	if iterationDuration < distributedIterationDuration {
		return iterationDuration, rateFn
	}

	rate := 0
	accRate := 0.0
	remainingSteps := 0
	tickSteps := int(iterationDuration.Milliseconds() / distributedIterationDuration.Milliseconds())

	distributedRateFn := func(time time.Time) int {
		if remainingSteps == 0 {
			rate = rateFn(time)
			accRate = 0.0
			remainingSteps = tickSteps
		}

		accRate += float64(rate) / float64(tickSteps)
		accRate = math.Round(accRate*10_000_000) / 10_000_000
		remainingSteps--

		if accRate < 1 {
			return 0
		}

		roundedAccRate := int(accRate)
		accRate -= float64(roundedAccRate)

		return roundedAccRate
	}

	return distributedIterationDuration, distributedRateFn
}
