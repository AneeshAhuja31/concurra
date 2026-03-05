package policy

import "time"

type errorAware struct {
	prevQueueDepth int
}

func (p *errorAware) Name() string { return "error-aware" }

func (p *errorAware) Evaluate(snapshot Snapshot, currentLimit int) Decision {
	defer func() { p.prevQueueDepth = snapshot.QueueDepth }()

	//priority1: error rate > 5% -> decrease 35%
	if snapshot.ErrorRate > 0.05 {
		mag := currentLimit * 35 / 100
		if mag < 1 {
			mag = 1
		}
		return Decision{Direction: Decrease, Magnitude: mag, Reason: "error rate exceeds 5%"}
	}
	//priority2: error rate > 1% -> decrease 15%
	if snapshot.ErrorRate > 0.01 {
		mag := currentLimit * 15 / 100
		if mag < 1 {
			mag = 1
		}
		return Decision{Direction: Decrease, Magnitude: mag, Reason: "error rate exceeds 1%"}
	}
	//priority3: low error rate + low latency + queue shrinking -> increase 5%
	queueShrinking := snapshot.QueueDepth < p.prevQueueDepth
	if snapshot.ErrorRate < 0.005 && snapshot.LatencyP99 < 200*time.Millisecond && queueShrinking {
		mag := currentLimit * 5/100
		if mag <1{
			mag = 1
		}
		return Decision{Direction: Increase, Magnitude: mag, Reason: "low error rate, low latency, queue shrinking"}
	}
	return Decision{Direction: Hold, Reason: "error rate within tolerance" }
}