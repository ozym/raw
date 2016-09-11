package raw

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Reader interface {
	Read(io.Reader) ([]Reading, error)
}

func Read(r io.Reader, rd Reader) ([]Reading, error) {
	return rd.Read(r)
}

func ReadFile(path string, rd Reader) ([]Reading, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return rd.Read(f)
}

type Writer interface {
	Write(io.Writer, []Reading) error
}

func Write(w io.Writer, wr Writer, readings []Reading) error {
	return wr.Write(w, readings)
}

func WriteFile(path string, wr Writer, readings []Reading) error {

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	defer os.Chmod(path, 0644)

	f, err := ioutil.TempFile(filepath.Dir(path), ".xxxx")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	if err := wr.Write(f, readings); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	if err := os.Rename(f.Name(), path); err != nil {
		return err
	}

	return nil
}

type ReadWriter interface {
	Reader
	Writer
}

func ReadWriteFile(path string, rw ReadWriter, readings []Reading) error {
	if _, err := os.Stat(path); err != nil {
		return WriteFile(path, rw, readings)
	}

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return WriteFile(path, rw, readings)
	}

	existing, err := Read(bytes.NewBuffer(raw), rw)
	if err != nil {
		return WriteFile(path, rw, readings)
	}

	obs := Merge(existing, readings)

	var buf bytes.Buffer
	if err := Write(&buf, rw, obs); err != nil {
		return WriteFile(path, rw, readings)
	}

	if !bytes.Equal(raw, buf.Bytes()) {
		return WriteFile(path, rw, readings)
	}

	return nil
}

func Store(dir string, rw ReadWriter, filename func(Reading) (string, error), readings []Reading) error {

	// map readings into files
	files := make(map[string][]Reading)
	for _, r := range readings {
		n, err := filename(r)
		if err != nil {
			return err
		}
		files[n] = append(files[n], r)
	}

	// update each file
	for k, rr := range files {
		if err := ReadWriteFile(filepath.Join(dir, k), rw, rr); err != nil {
			return err
		}
	}

	return nil
}
