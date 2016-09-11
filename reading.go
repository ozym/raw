package raw

import (
	"strconv"
	"strings"
	"time"
)

type Reading struct {
	Source string
	Epoch  time.Time
	Value  float64
}

func (r Reading) Less(reading Reading) bool {
	switch {
	case r.Source < reading.Source:
		return true
	case r.Source > reading.Source:
		return false
	default:
		return r.Epoch.Before(reading.Epoch)
	}
}

func (r Reading) Equal(reading Reading) bool {
	switch {
	case r.Less(reading):
		return false
	case reading.Less(r):
		return false
	default:
		return true
	}
}

func (r Reading) Date() string {
	b, err := r.Epoch.MarshalText()
	if err != nil {
		b, _ = time.Unix(0, 0).MarshalText()
	}
	return string(b)
}

func (r Reading) Key() string {
	return strings.Join([]string{r.Source, r.Date()}, ":")
}

func (r Reading) String() string {
	return strings.Join([]string{r.Source, r.Date(), strconv.FormatFloat(r.Value, 'f', -1, 64)}, " ")
}
