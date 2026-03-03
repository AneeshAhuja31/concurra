package ring

//Ring is a fixed size circular buffer. when full, new pushes overwrite the oldest entry.
type Ring[T any] struct{
	buf []T
	head int //next write position
	count int
	cap int
}

func New[T any](capacity int) *Ring[T] {
	if capacity < 1 {
		panic("ring capacity must be >= 1")
	}
	return &Ring[T]{
		buf: make([]T, capacity),
		cap: capacity,
	}
}

func (r *Ring[T]) Push(v T){
	r.buf[r.head] = v
	r.head = (r.head + 1) % r.cap
	if r.count < r.cap {
		r.count++
	}
}

func (r *Ring[T]) Len() int{
	return r.count
}

func (r *Ring[T]) Cap() int {
	return r.cap
}
// returns all valid entries in chronological order oldest to newest
func (r *Ring[T]) Slice() []T{ 
	if r.count == 0 {
		return []T{}
	}
	result := make([]T, r.count)
	start := (r.head - r.count) % r.cap
	for i := 0; i < r.count; i++ {
		result[i] = r.buf[(start+i)%r.cap]
	}
	return result
}

func (r *Ring[T]) Clear() {
	r.head=0
	r.count=0
}