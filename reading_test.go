package geomag

import (
	"io/ioutil"
	"math"
	"os"
	"testing"
	"time"
)

func TestReading_Scale(t *testing.T) {
	var readings = []Reading{
		{"", time.Now(), -1.0},
		{"", time.Now(), 0.0},
		{"", time.Now(), 1.0},
	}

	var scales = []float64{1.0, -1.0, 2.0, -2.0}

	for _, s := range scales {
		t.Logf("checking scale: %g", s)
		scaled := Readings(readings).Scale(s)
		if len(scaled) != len(readings) {
			t.Fatalf("invalid scaled length")
		}
		for i, r := range readings {
			if math.Abs(s*r.Value-scaled[i].Value) > 1.0e-09 {
				t.Errorf("invalid scaling %d: %s <= %s", i, r.String(), scaled[i].String())
			}
		}
	}

}

func TestReading_Merge(t *testing.T) {
	now := time.Now()

	var set1 = []Reading{
		{"a", now.Add(0 * time.Second), 0.0},
		{"a", now.Add(1 * time.Second), 1.0},
		{"a", now.Add(2 * time.Second), 2.0},
	}

	var set2 = []Reading{
		{"a", now.Add(2 * time.Second), 2.0},
		{"a", now.Add(3 * time.Second), 3.0},
		{"a", now.Add(4 * time.Second), 4.0},
		{"a", now.Add(5 * time.Second), 5.0},
	}

	var set3 = []Reading{
		{"b", now.Add(0 * time.Second), 0.0},
		{"b", now.Add(1 * time.Second), 1.0},
		{"b", now.Add(2 * time.Second), 2.0},
	}

	m1 := Readings(set1).Merge(set1)
	if len(m1) != len(set1) {
		t.Error("unable to merge set with itself")
	}

	m2 := Readings(set1).Merge(set2)
	if len(m2) != len(set1)+len(set2)-1 {
		t.Error("unable to merge sets")
	}
	m3 := Readings(set1).Merge(set3)
	if len(m3) != len(set1)+len(set3) {
		t.Error("unable to merge different sets")
	}

}

func TestReading_File(t *testing.T) {

	var tests = []struct {
		f string
		b [2]float64
		n int
	}{
		{
			"testdata/2016.215.04.NZ_APIM_50_LFZ.csv",
			[2]float64{-41221, -37449},
			3600,
		},
	}

	for _, x := range tests {
		t.Logf("checking file %s", x.f)

		f, err := os.OpenFile(x.f, os.O_RDONLY, 0)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()

		var r Readings
		if err := r.Read(f); err != nil {
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
		tf, err := ioutil.TempFile("", "")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tf.Name())
		if err := Readings(r).Write(tf, -1); err != nil {
			t.Fatal(err)
		}
		if err := tf.Close(); err != nil {
			t.Fatal(err)
		}
		if !equalFileContents(tf.Name(), x.f) {
			t.Errorf("unable to recover file exactly: %s", x.f)
		}

	}
}
