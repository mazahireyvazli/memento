package memento

import (
	"bytes"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"
)

const valueSize = 100 // in bytes

var config = &MementoConfig{
	ShardNum:       1 << 10,
	ShardCapHint:   1 << 16,
	EntryExpiresIn: time.Minute * 1,
}

func key(i int) string {
	return fmt.Sprintf("key-%010d", i)
}
func parallelKey(threadID int, counter int) string {
	return fmt.Sprintf("key-%04d-%06d", threadID, counter)
}
func value() []byte {
	return make([]byte, valueSize)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func BenchmarkSet(b *testing.B) {
	var memcache, _ = NewMemento[string](config)
	defer memcache.Close()

	for i := 0; i < b.N; i++ {
		memcache.Set(key(i), value())
	}
}

func BenchmarkGet(b *testing.B) {
	var memcache, _ = NewMemento[string](config)
	defer memcache.Close()

	for i := 0; i < b.N; i++ {
		memcache.Set(key(i), value())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v, ok := memcache.Get(key(i))

		if !ok || v == nil {
			b.Fatalf("couldn't find entry for provided key %s", key(i))
		}

		if !bytes.Equal(v, value()) {
			b.Fatalf("mismatch for provided key %s. expected %s, got %s", key(i), value(), v)
		}
	}
}

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

func BenchmarkParallelSet(b *testing.B) {
	var memcache, _ = NewMemento[string](config)
	defer memcache.Close()

	b.ResetTimer()
	var i int64
	b.RunParallel(func(pb *testing.PB) {
		id := int(atomic.AddInt64(&i, 1))
		c := id * b.N
		for ; pb.Next(); c++ {
			memcache.Set(parallelKey(id, c), value())
		}
	})
}

func BenchmarkParallelGet(b *testing.B) {
	var memcache, _ = NewMemento[string](config)
	defer memcache.Close()

	var i int64
	b.RunParallel(func(pb *testing.PB) {
		id := int(atomic.AddInt64(&i, 1))
		c := id * b.N
		for ; pb.Next(); c++ {
			memcache.Set(parallelKey(id, c), value())
		}
	})

	atomic.StoreInt64(&i, 0)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id := int(atomic.AddInt64(&i, 1))
		c := id * b.N

		for ; pb.Next(); c++ {
			memcache.Get(parallelKey(id, c))
		}
	})
}

func BenchmarkParallelDelete(b *testing.B) {
	var memcache, _ = NewMemento[string](config)
	defer memcache.Close()

	var i int64
	b.RunParallel(func(pb *testing.PB) {
		id := int(atomic.AddInt64(&i, 1))
		c := id * b.N

		for ; pb.Next(); c++ {
			memcache.Set(parallelKey(id, c), value())
		}
	})
	var lengthBeforeDelete = memcache.Length()

	atomic.StoreInt64(&i, 0)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id := int(atomic.AddInt64(&i, 1))
		c := id * b.N

		for ; pb.Next(); c++ {
			memcache.Delete(parallelKey(id, c))
		}
	})

	var lengthAfterDelete = memcache.Length()

	if !(lengthAfterDelete < lengthBeforeDelete) {
		b.Fatal("most items should've been deleted")
	}
}
