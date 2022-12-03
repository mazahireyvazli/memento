package memento

import (
	"bytes"
	"strconv"
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
	return "key-" + strconv.Itoa(i)
}
func parallelKey(threadID int, counter int) string {
	return "key-" + strconv.Itoa(threadID) + "-" + strconv.Itoa(counter)
}
func value() []byte {
	return make([]byte, valueSize)
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
