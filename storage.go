package geomag

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"text/template"
	"time"
)

const tmpReadingPrefix = ".xxxx"

func fileContents(f string) ([]byte, error) {
	if _, err := os.Stat(f); err != nil {
		return nil, err
	}
	return ioutil.ReadFile(f)
}

func equalFileContents(f1, f2 string) bool {
	x1, err := fileContents(f1)
	if err != nil {
		return false
	}
	x2, err := fileContents(f2)
	if err != nil {
		return false
	}
	return bytes.Equal(x1, x2)
}

type Storage struct {
	Template     *template.Template
	DecimalPlace int
}

func NewStorage(tmpl string, dp int) (*Storage, error) {
	t, err := template.New("readings").Funcs(template.FuncMap{
		"Year": func(t time.Time) string {
			return t.Format("2006")
		},
		"Month": func(t time.Time) string {
			return t.Format("01")
		},
		"Day": func(t time.Time) string {
			return t.Format("02")
		},
		"Doy": func(t time.Time) string {
			return fmt.Sprintf("%03d", t.YearDay())
		},
		"Hour": func(t time.Time) string {
			return t.Format("15")
		},
		"Minute": func(t time.Time) string {
			return t.Format("04")
		},
		"Second": func(t time.Time) string {
			return t.Format("05")
		},
	}).Parse(tmpl)
	if err != nil {
		return nil, err
	}

	return &Storage{
		Template:     t,
		DecimalPlace: dp,
	}, nil
}

func (s Storage) Store(dir string, readings []Reading) error {

	// map readings into files
	files := make(map[string][]Reading)
	for _, v := range readings {
		b := new(bytes.Buffer)
		if err := s.Template.Execute(b, v); err != nil {
			return err
		}
		files[b.String()] = append(files[b.String()], v)
	}

	// update each file
	for k, v := range files {
		path := filepath.Join(dir, k)

		m := Readings([]Reading{}).Merge(v)

		// import any existing readings
		if _, err := os.Stat(path); err == nil {
			if r, err := s.ReadFile(path); err != nil {
				m = Readings(m).Merge(r)
			}
		}

		// get them in order
		sort.Sort(Readings(m))

		// write out the readings if they're different
		if err := s.WriteFile(path, m); err != nil {
			return err
		}
	}

	return nil
}

func (s Storage) ReadFile(path string) ([]Reading, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var r Readings
	if err := r.Read(f); err != nil {
		return nil, err
	}

	return r, nil
}

func (s Storage) WriteFile(path string, readings []Reading) error {

	// write to a temporary file first and then rename it
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	defer os.Chmod(path, 0644)

	f, err := ioutil.TempFile(filepath.Dir(path), tmpReadingPrefix)
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if err := Readings(readings).Write(f, s.DecimalPlace); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if !equalFileContents(f.Name(), path) {
		if err := os.Rename(f.Name(), path); err != nil {
			return err
		}
	}

	return nil
}
