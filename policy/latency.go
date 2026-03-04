package policy

import "time"

type latencySensitive struct {
	consecutiveHighLatency int
	ticksSinceLastDecrease int
}

func (p *latencySensitive) Name() string { return "latency-sensitive" }

func (p *latencySensitive) Evaluate(snapshot Snapshot, currentLimit int) Decision {
	//priority1: Emergency - P99 > 350ms -> decrease 25%
	if snapshot.LatencyP99 > 350*time.Millisecond {
		p.consecutiveHighLatency++
		p.ticksSinceLastDecrease = 0
		mag := currentLimit * 25 / 100
		if mag < 1 {
			mag = 1
		}
		return Decision{Direction: Decrease, Magnitude: mag, Reason: "P99 latency exceeds 350ms - emergency reduction"}
	}

	// priority 2: P99 > 200ms for 2+ consecutive evaluations -> decrease 10%
	if snapshot.LatencyP99 > 200*time.Millisecond {
		p.consecutiveHighLatency++
		if p.consecutiveHighLatency >= 2 {
			p.ticksSinceLastDecrease = 0
			mag := currentLimit * 10 / 100
			if mag < 1 {
				mag = 1
			}
			return Decision{Direction: Decrease, Magnitude: mag, Reason: "P99 latency above 200ms for 2+ consecutive ticks"}
		}
		//first tick above 200ms — hold and wait for confirmation.
		p.ticksSinceLastDecrease++
		return Decision{Direction: Hold, Reason: "P99 elevated, awaiting confirmation"}
	}

	//latency is acceptable — reset consecutive high latency counter.
	p.consecutiveHighLatency = 0
	p.ticksSinceLastDecrease++

	//priority 3: Queue building with low latency -> increase 5%
	if snapshot.QueueDepth > 10 && snapshot.LatencyP99 < 150*time.Millisecond {
		mag := currentLimit * 5 / 100
		if mag < 1 {
			mag = 1
		}
		return Decision{Direction: Increase, Magnitude: mag, Reason: "queue building with low latency"}
	}

	// priority 4: System stable for 3+ minutes (~360 ticks at 500ms interval, but we count evaluations).
	// We approximate "3 minutes" as 360 ticks assuming default 500ms interval.
	// More practically, we track ticks since last decrease.
	if p.ticksSinceLastDecrease >= 360 {
		mag := currentLimit * 5 / 100
		if mag < 1 {
			mag = 1
		}
		return Decision{Direction: Increase, Magnitude: mag, Reason: "system stable for 3+ minutes — cautious increase"}
	}

	return Decision{Direction: Hold, Reason: "system stable"}
}
