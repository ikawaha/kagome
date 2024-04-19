package mem_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/ikawaha/kagome/v2/tokenizer/lattice/mem"
)

func ExamplePool() {
	type Foo struct {
		Bar string
	}

	// newFoo is a constructor of Foo type.
	newFoo := func() *Foo {
		return &Foo{
			Bar: "let the foo begin",
		}
	}

	// Create a new memory pool of Foo type.
	// If the memory pool is empty, it creates a new instance of Foo using newFoo.
	bufPool := mem.NewPool[Foo](newFoo)

	// Retrieve a Foo instance from the memory pool and print the current value
	// of the Bar field.
	a := bufPool.Get()
	fmt.Println(a.Bar)

	// Set the Bar field then put it back to the memory pool.
	a.Bar = "buz"
	bufPool.Put(a)

	// Same as above but set a different value to the Bar field.
	//
	// This will overwrite the previous value. But note that this will not allocate
	// new memory and is safe for use by multiple goroutines simultaneously.
	// See the benchmark in the same test file.
	b := bufPool.Get()
	b.Bar = "qux"
	bufPool.Put(b)

	// Retrieve a Foo instance from the memory pool and print the current value
	// of the Bar field.
	c := bufPool.Get()
	fmt.Println(c.Bar)
	// Output:
	// let the foo begin
	// qux
}

// To benchmark run:
//
//	go test -benchmem -bench=Benchmark -count 5 ./tokenizer/lattice/mem
func BenchmarkPool(b *testing.B) {
	type Foo struct {
		Bar int
	}

	bufPool := mem.NewPool[Foo](func() *Foo {
		return new(Foo)
	})

	b.ResetTimer()

	// Spawn 3 goroutines to get and put the Foo instance from the memory pool
	// b.N times each.
	//
	// Note that mem.Pool is safe for use by multiple goroutines simultaneously
	// and no new memory will be allocated.
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)

		go func(max int) {
			defer wg.Done()

			for i := 0; i < max; i++ {
				a := bufPool.Get() // get
				a.Bar += i         // increment
				bufPool.Put(a)     // put
			}
		}(b.N)
	}

	wg.Wait()
}

// To fuzz test run:
//
//	go test -fuzz=Fuzz -fuzztime 1m ./tokenizer/lattice/mem/...
func FuzzPool(f *testing.F) {
	type Foo struct {
		Bar string
	}

	bufPool := mem.NewPool[Foo](func() *Foo {
		return new(Foo)
	})

	// Corpus variants/patterns to fuzz test.
	f.Add("")
	f.Add(" ")
	f.Add("0123456789")
	f.Add("!@#$%^&*()_+")
	f.Add("short")
	f.Add("短")
	f.Add("this is a test with a long string")
	f.Add("これは、いささか長いテスト用の文字列です。")

	f.Fuzz(func(t *testing.T, s string) {
		a := bufPool.Get()
		a.Bar += s
		bufPool.Put(a)
	})
}
