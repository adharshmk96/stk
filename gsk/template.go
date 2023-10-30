package gsk

import (
	"bytes"
	"html/template"
)

type Tpl struct {
	TemplatePath string
	Variables    interface{}
}

type comboVariables struct {
	Var    interface{}
	Config map[string]interface{}
}

func (t *Tpl) Render(configVars map[string]interface{}) ([]byte, error) {
	tmpl, err := template.ParseFiles(t.TemplatePath)
	if err != nil {
		return nil, err
	}

	comboVars := &comboVariables{
		Var:    t.Variables,
		Config: configVars,
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, comboVars)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
