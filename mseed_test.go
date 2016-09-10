package geomag

import (
	"math"
	"testing"
)

func TestMSeed_Scale(t *testing.T) {
	var valid = []float64{-1.0, 1.0}

	for _, v := range valid {
		t.Logf("checking valid scale %g", v)
		if _, err := NewMSeed(v); err != nil {
			t.Errorf("scale %g should be valid: %v", v, err)
		}
	}

	var invalid = []float64{0.0}

	for _, v := range invalid {
		t.Logf("checking invalid scale %g", v)
		if _, err := NewMSeed(v); err == nil {
			t.Errorf("scale %g should be invalid", v)
		}
	}

}

func TestMSeed_File(t *testing.T) {

	var tests = []struct {
		f string
		s float64
		b [2]float64
		n int
	}{
		{
			"testdata/NZ.APIM.50.LFZ.D.2016.215",
			1.0,
			[2]float64{-55042, -29273},
			83955,
		},
	}

	for _, x := range tests {
		t.Logf("checking file %s", x.f)

		m, err := NewMSeed(x.s)
		if err != nil {
			t.Fatal(err)
		}

		r, err := m.ReadFile(x.f)
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

	}
}
