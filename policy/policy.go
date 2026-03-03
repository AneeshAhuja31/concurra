//package polict defines the Policy interface and built in concurrency adjustment policies.
//a policy evaluates system metrics and decides whether to increase, decrease or hold the current concurrency limits
//4 built in policies are provided:
// LatencySensitive, ErrorAware, QueueFirst, and StabilityFirst
package policy

import "time"

type Direction int

const (
	Hold Direction = iota
	Increase
	Decrease
)

func (d Direction) String() string{
	switch d {
	case Hold:
		return "hold"
	case Increase:
		return "increase"
	case Decrease:
		return "decrease"
	default:
		return "unknown"
	}
}

//Snapshot contains a point-in-time view of system metrics used by policies for decision making
type Snapshot struct {
	LatencyP50 time.Duration // 50th percentile observed latency
	LatencyP95 time.Duration
	LatencyP99 time.Duration
	ErrorRate float64 //fraction of requests that failed [0.0,1.0]
	QueueDepth int //number of tasks waiting for a concurrent slot
	ActiveSlots int 
	WindowSize int
}
// Decision describes a concurrency adjustment recommended by a Policy
type Decision struct{
	Direction Direction
	Magnitude int
	Reason string
}

type Policy interface{
	Evaluate(snapshot Snapshot, currentLimit int) Decision
	Name() string 
}


