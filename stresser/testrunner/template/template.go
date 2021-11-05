package template

import (
	"strings"
	"text/template"
)

func Format(fmt string, data interface{}) (string, error) {
	tmpl := template.New("")
	if _, err := tmpl.Parse(fmt); err != nil {
		return "", err
	}
	b := new(strings.Builder)
	if err := tmpl.Execute(b, data); err != nil {
		return "", err
	}
	return b.String(), nil
}
