package memento

import (
	"sync/atomic"
	"testing"
)

func BenchmarkParallelSet(b *testing.B) {
	var memcache, _ = NewMemento[string](config)
	defer memcache.Close()

	b.ResetTimer()
	var i int64
	b.RunParallel(func(pb *testing.PB) {
		id := int(atomic.AddInt64(&i, 1) - 1)

		for c := 0; pb.Next(); c++ {
			memcache.Set(parallelKey(id, c), value())
		}
	})
}

func BenchmarkParallelGet(b *testing.B) {
	var memcache, _ = NewMemento[string](config)
	defer memcache.Close()

	var i int64
	b.RunParallel(func(pb *testing.PB) {
		id := int(atomic.AddInt64(&i, 1) - 1)

		for c := 0; pb.Next(); c++ {
			memcache.Set(parallelKey(id, c), value())
		}
	})

	atomic.StoreInt64(&i, 0)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id := int(atomic.AddInt64(&i, 1) - 1)

		for c := 0; pb.Next(); c++ {
			memcache.Get(parallelKey(id, c))
		}
	})
}
