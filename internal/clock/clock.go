//package clock provide a mockable time interface for ttesting
package clock

import (
	"sync"
	"time"
)

type Clock interface{
	Now() time.Time
	After(d time.Time) <- chan time.Time
	Since(t time.Time) time.Duration
}

type Real struct{}

func (Real) Now() time.Time { return time.Now() }

func (Real) After(d time.Duration) <-chan time.Time { return time.After(d) }

func (Real) Since(t time.Time) time.Duration { return time.Since(t)}

type Mock struct{
	mu sync.Mutex
	now time.Time
	waiters []mockWaiter //list of timers waiting for future deadlines
}

type mockWaiter struct {
	deadline time.Time
	ch chan time.Time
}

func NewMock(initial time.Time) *Mock{
	return &Mock{now: initial}
}

func (m *Mock) Now() time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.now
}
// After returns a channel that fires when the mock time is advanced past the deadline
func (m *Mock) After(d time.Duration) <-chan time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	ch := make(chan time.Time, 1)
	deadline := m.now.Add(d)
	if !m.now.Before(deadline){ //if deadline is in past send it immediately
		ch <- m.now
		return ch
	}

	m.waiters  = append(m.waiters, mockWaiter{deadline: deadline, ch: ch})
	return ch
}
// returns mock time elapsed since t
func (m *Mock) Since(t time.Time) time.Duration{
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.now.Sub(t)
}

func (m *Mock) Advance(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.now = m.now.Add(d)
	m.fireWaiters()
}

func (m *Mock) fireWaiters(){
	remaining := m.waiters[:0]
	for _,w := range m.waiters{
		if !m.now.Before(w.deadline){
			w.ch <- m.now
		} else {
			remaining = append(remaining, w)
		}
	}
	m.waiters = remaining
}

func (m *Mock) Set(t time.Time){
	m.mu.Lock()
	defer m.mu.Unlock()
	m.now = t
	m.fireWaiters()
}