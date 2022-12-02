package memento

import (
	"sync"
	"testing"
	"time"
)

func BenchmarkRaceSetAndGet(b *testing.B) {
	var iter_num = b.N

	var memcache, _ = NewMemento[string](config)
	defer memcache.Close()

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		for c := 0; c < iter_num; c++ {
			memcache.Set(parallelKey(c, c), value())
		}

		wg.Done()
	}()

	go func() {
		time.Sleep(time.Millisecond * (time.Duration(iter_num / 5000)))

		var misses int

		for c := 0; c < iter_num; c++ {
			v, ok := memcache.Get(parallelKey(c, c))

			if v == nil || !ok {
				misses++
			}
		}

		b.Log("total misses", misses)

		wg.Done()
	}()

	go func() {
		time.Sleep(time.Millisecond * (time.Duration(iter_num / 5000)))
		b.Log("total items in cache before delete", memcache.Length())

		for c := 0; c < iter_num; c++ {
			memcache.Delete(parallelKey(c, c))
		}

		b.Log("total items in cache after delete", memcache.Length())

		wg.Done()
	}()

	wg.Wait()

}
