package raw

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"time"
)

const (
	csvEpochIndex int = iota
	csvSourceIndex
	csvValueIndex
	csvLastIndex
)

type Csv struct {
	DecimalPlace *int
}

func NewCsv(dp int) *Csv {
	return &Csv{
		DecimalPlace: &dp,
	}
}

func (c Csv) Read(rd io.Reader) ([]Reading, error) {

	var readings []Reading
	data, err := csv.NewReader(rd).ReadAll()
	if err != nil {
		return nil, err
	}
	for n, d := range data {
		if len(d) != csvLastIndex {
			return nil, fmt.Errorf("line %d: invalid sample element length: %d", n, len(d))
		}

		var t time.Time
		if err := t.UnmarshalText([]byte(d[csvEpochIndex])); err != nil {
			return nil, fmt.Errorf("line %d: invalid sample time: %v", n, err)
		}

		v, err := strconv.ParseFloat(d[csvValueIndex], 64)
		if err != nil {
			return nil, fmt.Errorf("line %d: invalid sample float: %v", n, err)
		}

		readings = append(readings, Reading{
			Source: d[csvSourceIndex],
			Epoch:  t,
			Value:  v,
		})
	}

	return readings, nil
}

func (c Csv) Write(wr io.Writer, rr []Reading) error {
	var dp int

	switch {
	case c.DecimalPlace != nil:
		dp = *c.DecimalPlace
	default:
		dp = -1
	}

	data := [][]string{}
	for _, r := range rr {
		b, err := r.Epoch.MarshalText()
		if err != nil {
			return err
		}
		data = append(data, []string{string(b), r.Source, strconv.FormatFloat(r.Value, 'f', dp, 64)})
	}
	if err := csv.NewWriter(wr).WriteAll(data); err != nil {
		return err
	}

	return nil
}
