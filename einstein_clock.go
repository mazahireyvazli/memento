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

func (c *EinsteinClock) updateTime(t time.Time) {
	c.seconds = uint64(t.Unix())
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
		seconds: uint64(time.Now().Unix()),
		donech:  make(chan struct{}),
	}

	clock.timeUpdater()

	return clock
}
