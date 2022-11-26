package memento

import "time"

type EinsteinClock struct {
	seconds uint64
	donech  chan struct{}
}

func (c *EinsteinClock) Close() {
	close(c.donech)
}

func (c *EinsteinClock) Seconds() uint64 {
	return c.seconds
}

func (c *EinsteinClock) UpdateTime() {
	c.seconds = uint64(time.Now().Unix())
}

func (c *EinsteinClock) TimeUpdater() {
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c.UpdateTime()
			case <-c.donech:
				return
			}
		}
	}()
}

func NewClock() (clock *EinsteinClock) {
	clock = &EinsteinClock{
		seconds: uint64(time.Now().Unix()),
		donech:  make(chan struct{}),
	}

	clock.TimeUpdater()

	return clock
}
