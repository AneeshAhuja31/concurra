package policy


type latencySensitive struct {
	consecutiveHighLatency int
	ticksSinceLastDecrease int
}

func (p *latencySensitive) Name() string { return "latency-sensitive" }

func (p *latencySensitive) Evaluate(snapshot Snapshot, currentLimit int) Decision {

}