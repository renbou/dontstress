package util

import (
	"strings"
	"text/template"
)

type KV map[string]string

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
