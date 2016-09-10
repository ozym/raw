package geomag

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

const (
	networkIndex int = iota
	stationIndex
	locationIndex
	channelIndex
)

const (
	epochIndex int = iota
	sourceIndex
	valueIndex
	lastIndex
)

type Reading struct {
	Source string
	Epoch  time.Time
	Value  float64
}

func (r Reading) source(part int) string {
	if parts := strings.Split(r.Source, "_"); len(parts) > part {
		return parts[part]
	}
	return ""
}

func (r Reading) Network() string {
	return r.source(networkIndex)
}

func (r Reading) Station() string {
	return r.source(stationIndex)
}

func (r Reading) Location() string {
	return r.source(locationIndex)
}

func (r Reading) Channel() string {
	return r.source(channelIndex)
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

func (r Reading) Encode(dp int) []string {
	return []string{r.Date(), r.Source, strconv.FormatFloat(r.Value, 'f', dp, 64)}
}

func (r *Reading) Decode(data []string) error {

	if len(data) != lastIndex {
		return fmt.Errorf("invalid sample element length: %d", len(data))
	}

	var t time.Time
	if err := t.UnmarshalText([]byte(data[epochIndex])); err != nil {
		return err
	}

	v, err := strconv.ParseFloat(data[valueIndex], 64)
	if err != nil {
		return err
	}

	*r = Reading{
		Source: data[sourceIndex],
		Epoch:  t,
		Value:  v,
	}

	return nil
}

func (r Reading) String() string {
	return strings.Join(r.Encode(-1), " ")
}

func (r Reading) Key() string {
	return strings.Join([]string{r.Source, r.Date()}, ":")
}

type Readings []Reading

func (r Readings) Len() int           { return len(r) }
func (r Readings) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r Readings) Less(i, j int) bool { return r[i].Less(r[j]) }

func (r *Readings) Read(rd io.Reader) error {

	data, err := csv.NewReader(rd).ReadAll()
	if err != nil {
		return err
	}
	for _, d := range data {
		var v Reading
		if err := v.Decode(d); err != nil {
			return err
		}
		*r = append(*r, v)
	}

	return nil
}

func (r Readings) Write(wr io.Writer, dp int) error {

	data := [][]string{}
	for _, d := range r {
		data = append(data, d.Encode(dp))
	}
	if err := csv.NewWriter(wr).WriteAll(data); err != nil {
		return err
	}

	return nil
}

func (r Readings) Merge(readings []Reading) []Reading {
	var list = make(map[string]Reading)

	for _, v := range append(r, readings...) {
		list[v.Key()] = v
	}

	var joined []Reading
	for _, v := range list {
		joined = append(joined, v)
	}

	return joined
}

func (r Readings) Scale(scale float64) []Reading {
	var scaled []Reading
	for _, v := range r {
		scaled = append(scaled, Reading{
			Source: v.Source,
			Epoch:  v.Epoch,
			Value:  scale * v.Value,
		})
	}
	return scaled
}
