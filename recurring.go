package recurring

import "time"

// Recurring holds a channel that delivers 'ticks' of a clock at a given time of day (UTC).
type Recurring struct {
	C      <-chan time.Time // The channel on which the ticks are delivered
	c      chan time.Time
	ticker *time.Ticker
	quit   chan interface{}
}

func deadline(start time.Time, hour, min, sec, nsec int) time.Duration {
	now := start.UTC()
	nextTick := time.Date(now.Year(), now.Month(), now.Day(), hour, min, sec, nsec, time.UTC)
	if nextTick.Before(now) || nextTick.Equal(now) {
		nextTick = nextTick.Add(24 * time.Hour)
	}
	return nextTick.Sub(now)
}

// New returns a new Recurring containing a channel that will send the time at a given time of day (UTC)
// specified by the arguments. It will skip ticks if a receiver is slow enough (ie, more than 24 hours).
// Stop the Recurring to release associated resources.
func New(hour, min, sec, nsec int) *Recurring {
	now := time.Now().UTC()
	first := deadline(now, hour, min, sec, nsec)
	r := &Recurring{
		ticker: time.NewTicker(first),
		quit:   make(chan interface{}),
		c:      make(chan time.Time),
	}
	r.C = r.c
	go func() {
		for {
			select {
			case t := <-r.ticker.C:
				r.ticker.Stop()
				r.c <- t
				now := time.Now().UTC()
				r.ticker = time.NewTicker(deadline(now, hour, min, sec, nsec))
			case <-r.quit:
				r.ticker.Stop()
				return
			}
		}
	}()
	return r
}

// Stop turns off a Recurring. After Stop, no more ticks will be sent.
// Stop does not close the channel, to prevent a read from the channel succeeding incorrectly.
func (r *Recurring) Stop() {
	close(r.quit)
}
