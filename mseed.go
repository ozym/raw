package geomag

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/GeoNet/mseed"
)

type MSeed struct {
	Scale float64
}

func NewMSeed(scale float64) (*MSeed, error) {
	if scale == 0.0 {
		return nil, fmt.Errorf("scale cannot be zero")
	}
	return &MSeed{
		Scale: scale,
	}, nil
}

func (m MSeed) DecodeBuffer(buf []byte) ([]Reading, error) {
	var readings []Reading

	msr := mseed.NewMSRecord()
	defer mseed.FreeMSRecord(msr)

	msr.Unpack(buf, 512, 1, 0)
	samples, err := msr.DataSamples()
	if err != nil {
		return nil, err
	}

	sps := msr.Samprate()
	if len(samples) > 0 && sps > 0.0 {
		dt := time.Duration(float64(time.Second) / float64(sps))
		for n, s := range samples {
			readings = append(readings, Reading{
				Source: msr.SrcName(0),
				Epoch:  msr.Starttime().Add(time.Duration(n) * dt),
				Value:  m.Scale * float64(s),
			})
		}
	}

	return readings, nil
}

func (m MSeed) ReadStream(rd io.Reader) ([]Reading, error) {

	var readings []Reading

	// make space for miniseed blocks
	msr := mseed.NewMSRecord()
	defer mseed.FreeMSRecord(msr)

	buf := make([]byte, 512)
	for {
		if n, _ := io.ReadFull(rd, buf); n != len(buf) {
			break
		}
		r, err := m.DecodeBuffer(buf)
		if err != nil {
			return nil, err
		}

		readings = append(readings, r...)
	}

	return readings, nil
}

func (m MSeed) ReadFile(path string) ([]Reading, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := m.ReadStream(f)
	if err != nil {
		return nil, err
	}

	return r, nil
}
