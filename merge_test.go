package raw

import (
	"testing"
	"time"
)

func TestReading_Merge(t *testing.T) {
	now := time.Now()

	min := func(a, b int) int {
		if a < b {
			return a
		}
		return b
	}

	var tests = []struct {
		a, b, c []Reading
	}{
		// merge itself
		{
			[]Reading{
				{"a", now.Add(0 * time.Second), 0.0},
				{"a", now.Add(1 * time.Second), 1.0},
				{"a", now.Add(2 * time.Second), 2.0},
			},
			[]Reading{
				{"a", now.Add(0 * time.Second), 0.0},
				{"a", now.Add(1 * time.Second), 1.0},
				{"a", now.Add(2 * time.Second), 2.0},
			},
			[]Reading{
				{"a", now.Add(0 * time.Second), 0.0},
				{"a", now.Add(1 * time.Second), 1.0},
				{"a", now.Add(2 * time.Second), 2.0},
			},
		},
		// merge overlap
		{
			[]Reading{
				{"a", now.Add(0 * time.Second), 0.0},
				{"a", now.Add(1 * time.Second), 1.0},
				{"a", now.Add(2 * time.Second), 2.0},
			},
			[]Reading{
				{"a", now.Add(2 * time.Second), 2.0},
				{"a", now.Add(3 * time.Second), 3.0},
				{"a", now.Add(4 * time.Second), 4.0},
			},
			[]Reading{
				{"a", now.Add(0 * time.Second), 0.0},
				{"a", now.Add(1 * time.Second), 1.0},
				{"a", now.Add(2 * time.Second), 2.0},
				{"a", now.Add(3 * time.Second), 3.0},
				{"a", now.Add(4 * time.Second), 4.0},
			},
		},
		// merge different
		{
			[]Reading{
				{"a", now.Add(0 * time.Second), 0.0},
				{"a", now.Add(1 * time.Second), 1.0},
				{"a", now.Add(2 * time.Second), 2.0},
			},
			[]Reading{
				{"b", now.Add(0 * time.Second), 0.0},
				{"b", now.Add(1 * time.Second), 1.0},
				{"b", now.Add(2 * time.Second), 2.0},
			},
			[]Reading{
				{"a", now.Add(0 * time.Second), 0.0},
				{"a", now.Add(1 * time.Second), 1.0},
				{"a", now.Add(2 * time.Second), 2.0},
				{"b", now.Add(0 * time.Second), 0.0},
				{"b", now.Add(1 * time.Second), 1.0},
				{"b", now.Add(2 * time.Second), 2.0},
			},
		},
		// merge sort
		{
			[]Reading{
				{"a", now.Add(1 * time.Second), 1.0},
				{"a", now.Add(0 * time.Second), 0.0},
				{"a", now.Add(2 * time.Second), 2.0},
			},
			[]Reading{
				{"a", now.Add(2 * time.Second), 2.0},
				{"a", now.Add(1 * time.Second), 1.0},
				{"a", now.Add(0 * time.Second), 0.0},
			},
			[]Reading{
				{"a", now.Add(0 * time.Second), 0.0},
				{"a", now.Add(1 * time.Second), 1.0},
				{"a", now.Add(2 * time.Second), 2.0},
			},
		},
		// merge values
		{
			[]Reading{
				{"a", now.Add(0 * time.Second), 0.0},
				{"a", now.Add(1 * time.Second), 1.0},
				{"a", now.Add(2 * time.Second), 2.0},
			},
			[]Reading{
				{"a", now.Add(0 * time.Second), 3.0},
				{"a", now.Add(1 * time.Second), 4.0},
				{"a", now.Add(2 * time.Second), 5.0},
			},
			[]Reading{
				{"a", now.Add(0 * time.Second), 3.0},
				{"a", now.Add(1 * time.Second), 4.0},
				{"a", now.Add(2 * time.Second), 5.0},
			},
		},
	}

	for n, x := range tests {
		m := Merge(x.a, x.b)
		if len(m) != len(x.c) {
			t.Errorf("unable to merge test set [%d]", n)
		}
		for i := 0; i < min(len(m), len(x.c)); i++ {
			if m[i].Key() != x.c[i].Key() {
				t.Error("unable to merge test set [%d]: reading %d", n, i)
			}
			if m[i].Value != x.c[i].Value {
				t.Error("unable to merge test set [%d]: value %d", n, i)
			}
		}
	}

}
