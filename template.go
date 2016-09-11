package raw

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"
)

type Template struct {
	*template.Template
}

func NewTemplate(tmpl string) (*Template, error) {
	t, err := template.New("readings").Funcs(template.FuncMap{
		"Network": func(r Reading) string {
			if parts := strings.Split(r.Source, "_"); len(parts) > 0 {
				return parts[0]
			}
			return ""
		},
		"Station": func(r Reading) string {
			if parts := strings.Split(r.Source, "_"); len(parts) > 1 {
				return parts[1]
			}
			return ""
		},
		"Location": func(r Reading) string {
			if parts := strings.Split(r.Source, "_"); len(parts) > 2 {
				return parts[2]
			}
			return ""
		},
		"Channel": func(r Reading) string {
			if parts := strings.Split(r.Source, "_"); len(parts) > 3 {
				return parts[3]
			}
			return ""
		},
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

	return &Template{
		Template: t,
	}, nil
}

func (t Template) Execute(r Reading) (string, error) {
	if t.Template == nil {
		return "", fmt.Errorf("no template given")
	}
	b := new(bytes.Buffer)
	if err := t.Template.Execute(b, r); err != nil {
		return "", err
	}
	return b.String(), nil
}
