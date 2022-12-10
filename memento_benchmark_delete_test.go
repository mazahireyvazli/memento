package memento

import (
	"sync/atomic"
	"testing"

	"github.com/mazahireyvazli/memento/utils"
)

func BenchmarkDelete(b *testing.B) {
	var memcache, _ = NewMemento[string](config)
	defer memcache.Close()

	for i := 0; i < b.N; i++ {
		memcache.Set(utils.SimpleKey(i), utils.SimpleValue())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		memcache.Delete(utils.SimpleKey(i))
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
			memcache.Set(utils.ParallelKey(id, c), utils.SimpleValue())
		}
	})
	var lengthBeforeDelete = memcache.Length()

	atomic.StoreInt64(&i, 0)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id := int(atomic.AddInt64(&i, 1) - 1)

		for c := 0; pb.Next(); c++ {
			memcache.Delete(utils.ParallelKey(id, c))
		}
	})

	if !(memcache.Length() < lengthBeforeDelete) {
		b.Fatal("most items should've been deleted")
	}
}
