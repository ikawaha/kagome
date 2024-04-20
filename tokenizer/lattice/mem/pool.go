package mem

import (
	"sync"
)

// Pool represents memory pool of T.
//
// It is suitable for managing temporary objects that can be individually saved
// and retrieved (mem.Pool.Put and mem.Pool.Get).
// Unlike variables or pointer variables, mem.Pool is safe for use by multiple
// goroutines and no new memory will be allocated.
//
// It is a wrapper of sync.Pool to support generics and easy to use. For the
// actual usage, see the example in the package documentation.
type Pool[T any] struct {
	internal *sync.Pool
}

// NewPool returns a memory pool of T. newF is a constructor of T, and it is called
// when the memory pool is empty.
func NewPool[T any](newF func() *T) Pool[T] {
	return Pool[T]{
		internal: &sync.Pool{
			New: func() any {
				return newF()
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
