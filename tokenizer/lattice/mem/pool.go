package mem

import (
	"sync"
)

// Pool represents memory pool of T.
type Pool[T any] struct {
	internal *sync.Pool
}

// NewPool returns a memory pool of T.
func NewPool[T any](f func() *T) Pool[T] {
	return Pool[T]{
		internal: &sync.Pool{
			New: func() any {
				return f()
			},
		},
	}
}

// Get gets instance of T from the memory pool.
func (p Pool[T]) Get() *T {
	return p.internal.Get().(*T)
}

// Put puts the instance to the memory pool.
func (p Pool[T]) Put(x *T) {
	p.internal.Put(x)
}
