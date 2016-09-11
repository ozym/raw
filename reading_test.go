package raw

import (
	"testing"
	"time"
)

func TestReading_Keys(t *testing.T) {
	at := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)

	var tests = []struct {
		r Reading
		k string
	}{
		{
			Reading{"a", at, 1.0},
			"a:2010-01-01T00:00:00Z",
		}, {
			Reading{"b", at, 1.0},
			"b:2010-01-01T00:00:00Z",
		},
		{
			Reading{"c", at.Add(time.Second), 1.0},
			"c:2010-01-01T00:00:01Z",
		},
		{
			Reading{"d", at.Add(time.Second), 1.0},
			"d:2010-01-01T00:00:01Z",
		},
	}

	for _, x := range tests {
		if x.r.Key() != x.k {
			t.Errorf("unable to match key: %s != %s", x.r.Key(), x.k)
		}
	}

}

func TestReading_Strings(t *testing.T) {
	at := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)

	var tests = []struct {
		r Reading
		s string
	}{
		{
			// merge itself
			Reading{"a", at, 1.0},
			"a 2010-01-01T00:00:00Z 1",
		}, {
			Reading{"b", at, 1.0},
			"b 2010-01-01T00:00:00Z 1",
		},
		{
			Reading{"c", at.Add(time.Second), 1.0},
			"c 2010-01-01T00:00:01Z 1",
		},
		{
			Reading{"d", at.Add(time.Second), 2.0},
			"d 2010-01-01T00:00:01Z 2",
		},
	}

	for _, x := range tests {
		if x.r.String() != x.s {
			t.Errorf("unable to match string: %s != %s", x.r.String(), x.s)
		}
	}

}

func TestReading_Less(t *testing.T) {
	at := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)

	var tests = []struct {
		a, b Reading
		l    bool
	}{
		{
			Reading{"a", at, 1.0},
			Reading{"b", at, 1.0},
			true,
		}, {
			Reading{"a", at, 1.0},
			Reading{"a", at, 1.0},
			false,
		}, {
			Reading{"a", at.Add(time.Second), 1.0},
			Reading{"a", at, 1.0},
			false,
		}, {
			Reading{"a", at, 1.0},
			Reading{"a", at.Add(time.Second), 1.0},
			true,
		},
	}

	for _, x := range tests {
		if x.a.Less(x.b) != x.l {
			if x.l {
				t.Errorf("reading %s is not less than %s", x.a.Key(), x.b.Key())
			} else {
				t.Errorf("reading %s is less than %s", x.a.Key(), x.b.Key())
			}
		}
	}

}

func TestReading_Equal(t *testing.T) {
	at := time.Date(2010, 1, 1, 0, 0, 0, 0, time.UTC)

	var tests = []struct {
		a, b Reading
		l    bool
	}{
		{
			Reading{"a", at, 1.0},
			Reading{"b", at, 1.0},
			true,
		}, {
			Reading{"a", at, 1.0},
			Reading{"a", at.Add(time.Second), 2.0},
			true,
		}, {
			Reading{"a", at, 1.0},
			Reading{"a", at, 1.0},
			false,
		},
		{
			Reading{"a", at.Add(time.Second), 2.0},
			Reading{"a", at, 1.0},
			false,
		},
		{
			Reading{"b", at.Add(time.Second), 2.0},
			Reading{"a", at, 1.0},
			false,
		},
	}

	for _, x := range tests {
		if x.a.Less(x.b) != x.l {
			if x.l {
				t.Errorf("reading %s is not equal to %s", x.a.Key(), x.b.Key())
			} else {
				t.Errorf("reading %s is equal to %s", x.a.Key(), x.b.Key())
			}
		}
	}

}
