package scan

import (
	"time"
)

type Opts struct {
	minSize int64
	minAge  time.Duration
	now     time.Time
}

func DefaultOpts() Opts {
	return Opts{100 * 1 << 20, 30 * 24 * time.Hour, time.Now()}
}

func (o Opts) SetMinSize(s int64) Opts {
	o.minSize = s
	return o
}

func (o Opts) SetMinAge(d time.Duration) Opts {
	o.minAge = d
	return o
}

func (o Opts) SetNow(t time.Time) Opts {
	o.now = t
	return o
}
