package memento

import (
	"bytes"
	"testing"
	"time"

	"github.com/mazahireyvazli/memento/utils"
)

var config = &MementoConfig{
	ShardNum:       1 << 10,
	ShardCapHint:   1 << 16,
	EntryExpiresIn: time.Minute * 1,
}

func BenchmarkMemento(b *testing.B) {
	var memcache, _ = NewMemento[string](config)
	defer memcache.Close()

	b.Run("Set", func(b *testing.B) {
		if utils.SkipDiscoveryRun(b) {
			return
		}

		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			memcache.Set(utils.SimpleKey(i), utils.SimpleValue())
		}
	})

	b.Run("Get", func(b *testing.B) {
		if utils.SkipDiscoveryRun(b) {
			return
		}

		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			v, ok := memcache.Get(utils.SimpleKey(i))

			if !ok || v == nil {
				b.Fatalf("couldn't find entry for provided key %s", utils.SimpleKey(i))
			}

			if !bytes.Equal(v, utils.SimpleValue()) {
				b.Fatalf("mismatch for provided key %s. expected %s, got %s", utils.SimpleKey(i), utils.SimpleValue(), v)
			}
		}
	})
}
