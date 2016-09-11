package raw

import (
	"bytes"
	"io/ioutil"
	"math"
	"testing"
)

func TestCsv_File(t *testing.T) {

	var tests = []struct {
		c ReadWriter
		f string
		b [2]float64
		n int
	}{
		{
			Csv{},
			"testdata/2016.215.04.NZ_APIM_50_LFZ.csv",
			[2]float64{-41221, -37449},
			3600,
		},
	}

	for _, x := range tests {
		t.Logf("checking file %s", x.f)

		raw, err := ioutil.ReadFile(x.f)
		if err != nil {
			t.Fatal(err)
		}

		r, err := Read(bytes.NewBuffer(raw), x.c)
		if err != nil {
			t.Fatal(err)
		}

		if len(r) != x.n {
			t.Errorf("invalid number or records read for %s, expected %d found %d", x.f, x.n, len(r))
		}

		var min, max float64
		for i, v := range r {
			if i == 0 || min > v.Value {
				min = v.Value
			}
			if i == 0 || max < v.Value {
				max = v.Value
			}
		}

		if len(r) > 0 {
			if math.Abs(min-x.b[0]) > 1.0e-9 {
				t.Errorf("invalid minimum record value for %s, expected %g found %g", x.f, x.b[0], min)
			}
			if math.Abs(max-x.b[1]) > 1.0e-9 {
				t.Errorf("invalid maximum record value for %s, expected %g found %g", x.f, x.b[1], max)
			}
		}

		var buf bytes.Buffer
		if err := Write(&buf, x.c, r); err != nil {
			t.Fatal(err)
		}

		if !bytes.Equal(raw, buf.Bytes()) {
			t.Error("encoded and decoded record data should be the same: %s", x.f)
		}

	}
}
