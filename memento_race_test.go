package memento

import (
	"sync"
	"testing"
	"time"

	"github.com/mazahireyvazli/memento/utils"
)

func TestRaceSetGetAndDelete(t *testing.T) {
	const iter_num = 1_500_000

	var memcache, _ = NewMemento[string](config)
	defer memcache.Close()

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		for c := 0; c < iter_num; c++ {
			memcache.Set(utils.ParallelKey(c, c), utils.SimpleValue())
		}

		wg.Done()
	}()

	go func() {
		time.Sleep(time.Millisecond * 550)

		var misses int

		for c := 0; c < iter_num; c++ {
			v, ok := memcache.Get(utils.ParallelKey(c, c))

			if v == nil || !ok {
				misses++
			}
		}

		t.Log("total misses", misses)

		wg.Done()
	}()

	go func() {
		time.Sleep(time.Millisecond * 750)
		t.Log("total items in cache before delete", memcache.Length())

		for c := 0; c < iter_num; c++ {
			memcache.Delete(utils.ParallelKey(c, c))
		}

		t.Log("total items in cache after delete", memcache.Length())

		wg.Done()
	}()

	wg.Wait()
}
