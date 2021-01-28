package retry

import "time"

const defaultInterval = 3 * time.Second

type Repeat struct {
	times    int
	interval time.Duration
}

func (r *Repeat) Do(f func() error) error {
	return r.DoWithBreak(func() (error, bool) {
		err := f()
		return err, (err == nil)
	})
}

func (r *Repeat) DoWithBreak(f func() (error, bool)) error {
	var err error
	for i := 0; i < r.times; i++ {
		e, stop := f()

		if e == nil || stop {
			return e
		}

		err = e

		// Don't sleep if it's the last try
		if i < r.times-1 {
			time.Sleep(r.interval)
		}
	}
	return err
}

func Times(times int) *Repeat {
	return &Repeat{
		times:    times,
		interval: defaultInterval,
	}
}

func (r *Repeat) Interval(interval time.Duration) *Repeat {
	r.interval = interval
	return r
}

type Duration struct {
	limit    time.Duration
	interval time.Duration
}

func (d *Duration) Do(f func() error) error {
	var err error
	endTime := time.Now().Add(d.limit)
	for time.Now().Before(endTime) {
		err = f()
		if err == nil {
			return nil
		}

		if time.Now().Add(d.interval).Before(endTime) {
			time.Sleep(d.interval)
		} else {
			break
		}
	}
	return err
}

func For(t time.Duration) *Duration {
	return &Duration{
		limit:    t,
		interval: defaultInterval,
	}
}

func (d *Duration) Interval(interval time.Duration) *Duration {
	d.interval = interval
	return d
}
