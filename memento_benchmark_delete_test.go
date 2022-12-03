package memento

import (
	"sync/atomic"
	"testing"
)

func BenchmarkDelete(b *testing.B) {
	var memcache, _ = NewMemento[string](config)
	defer memcache.Close()

	for i := 0; i < b.N; i++ {
		memcache.Set(key(i), value())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		memcache.Delete(key(i))
	}

	if memcache.Length() != 0 {
		b.Fatal("all items should've been deleted")
	}
}

func BenchmarkParallelDelete(b *testing.B) {
	var memcache, _ = NewMemento[string](config)
	defer memcache.Close()

	var i int64
	b.RunParallel(func(pb *testing.PB) {
		id := int(atomic.AddInt64(&i, 1) - 1)

		for c := 0; pb.Next(); c++ {
			memcache.Set(parallelKey(id, c), value())
		}
	})
	var lengthBeforeDelete = memcache.Length()

	atomic.StoreInt64(&i, 0)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id := int(atomic.AddInt64(&i, 1) - 1)

		for c := 0; pb.Next(); c++ {
			memcache.Delete(parallelKey(id, c))
		}
	})

	var lengthAfterDelete = memcache.Length()

	if !(lengthAfterDelete < lengthBeforeDelete) {
		b.Fatal("most items should've been deleted")
	}
}
