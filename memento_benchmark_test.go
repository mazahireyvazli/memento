package memento

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const maxID = 1_000
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

func BenchmarkParallelSet(b *testing.B) {
	var memcache, _ = NewMemento[string](config)
	defer memcache.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		id := rand.Intn(maxID)
		counter := 0
		for pb.Next() {
			memcache.Set(parallelKey(id, counter), value())
			counter = counter + 1
		}
	})
}

func BenchmarkParallelGet(b *testing.B) {
	var memcache, _ = NewMemento[string](config)
	defer memcache.Close()

	b.RunParallel(func(pb *testing.PB) {
		counter := 0
		for pb.Next() {
			memcache.Set(key(counter), value())
			counter = counter + 1
		}
	})

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		var counter int
		for pb.Next() {
			memcache.Get(key(counter))
			counter++
		}
	})
}
