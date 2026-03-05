package policy

import "time"

type stabilityFirst struct {
	consecutiveStable int
}

func (p *stabilityFirst) Name() string { return "stability-first" }

func (p *stabilityFirst) Evaluate(snapshot Snapshot, currentLimit int) Decision {
	//aggressive decrease on any signal breach (same thresholds as errorAware).

	//emergency: error rate > 5% -> decrease 35%
	if snapshot.ErrorRate > 0.05 {
		p.consecutiveStable = 0
		mag := currentLimit * 35 / 100
		if mag < 1 {
			mag = 1
		}
		return Decision{Direction: Decrease, Magnitude: mag, Reason: "error rate exceeds 5% — emergency throttle"}
	}

	//error rate > 1% -> decrease 15%
	if snapshot.ErrorRate > 0.01 {
		p.consecutiveStable = 0
		mag := currentLimit * 15 / 100
		if mag < 1 {
			mag = 1
		}
		return Decision{Direction: Decrease, Magnitude: mag, Reason: "error rate exceeds 1%"}
	}

	// p99 > 200ms -> decrease 10%
	if snapshot.LatencyP99 > 200*time.Millisecond {
		p.consecutiveStable = 0
		mag := currentLimit * 10 / 100
		if mag < 1 {
			mag = 1
		}
		return Decision{Direction: Decrease, Magnitude: mag, Reason: "P99 latency exceeds 200ms"}
	}

	// No breach: count as stable.
	p.consecutiveStable++

	// Only increase after 5+ consecutive stable evaluations.
	// never increase by more than 2 slots at a time.
	if p.consecutiveStable >= 5 {
		mag := currentLimit * 5 / 100
		if mag < 1 {
			mag = 1
		}
		if mag > 2 {
			mag = 2
		}
		return Decision{Direction: Increase, Magnitude: mag, Reason: "stable for 5+ evaluations: cautious increase"}
	}

	return Decision{Direction: Hold, Reason: "awaiting stability confirmation"}
}
