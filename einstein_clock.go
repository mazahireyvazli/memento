package memento

import (
	"sync/atomic"
	"time"
)

type EinsteinClock struct {
	seconds atomic.Uint64
	donech  chan struct{}
}

func (c *EinsteinClock) Close() {
	close(c.donech)
}

func (c *EinsteinClock) Seconds() uint64 {
	return c.seconds.Load()
}

func (c *EinsteinClock) updateTime(t time.Time) {
	c.seconds.Store(uint64(t.Unix()))
}

func (c *EinsteinClock) timeUpdater() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case t := <-ticker.C:
				c.updateTime(t)
			case <-c.donech:
				return
			}
		}
	}()
}

func NewClock() (clock *EinsteinClock) {
	clock = &EinsteinClock{
		donech: make(chan struct{}),
	}
	clock.updateTime(time.Now())
	clock.timeUpdater()

	return clock
}
