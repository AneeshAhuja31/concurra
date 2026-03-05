package policy

import "time"

type queueFirst struct{}

func (p *queueFirst) Name() string { return "queue-first" }

func (p *queueFirst) Evaluate(snapshot Snapshot, currentLimit int) Decision {
	// p1: error rate > 2% -> decrease 10% regardless of queue
	if snapshot.ErrorRate > 0.02 {
		mag := currentLimit * 10 / 100
		if mag < 1 {
			mag = 1
		}
		return Decision{Direction: Decrease, Magnitude: mag, Reason: "error rate exceeds 2%"}
	}

	// p2: p99 > 500ms -> decrease 15% (higher tolerance than latency sensitive)
	if snapshot.LatencyP99 > 500*time.Millisecond {
		mag := currentLimit * 15 / 100
		if mag < 1 {
			mag = 1
		}
		return Decision{Direction: Decrease, Magnitude: mag, Reason: "P99 latency exceeds 500ms"}
	}

	// p3: deep queue > 50 -> increase 20%
	if snapshot.QueueDepth > 50 {
		mag := currentLimit * 20 / 100
		if mag < 1 {
			mag = 1
		}
		return Decision{Direction: Increase, Magnitude: mag, Reason: "queue depth exceeds 50 - aggressive increase"}
	}

	// p4: Queue > 20 -> increase 10%
	if snapshot.QueueDepth > 20 {
		mag := currentLimit * 10 / 100
		if mag < 1 {
			mag = 1
		}
		return Decision{Direction: Increase, Magnitude: mag, Reason: "queue depth exceeds 20"}
	}

	return Decision{Direction: Hold, Reason: "queue acceptable"}
}
