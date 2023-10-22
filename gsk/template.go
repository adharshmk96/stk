package gsk

import (
	"bytes"
	"html/template"
)

type Tpl struct {
	TemplatePath string
	Variables    interface{}
}

func (t *Tpl) Render() ([]byte, error) {
	tmpl, err := template.ParseFiles(t.TemplatePath)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, t.Variables)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
