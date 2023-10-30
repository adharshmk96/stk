package gsk_test

import (
	"testing"

	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/testutils"
	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	_, removeDir := testutils.SetupTempDirectory(t)
	defer removeDir()

	// Create test template files
	testutils.WriteFile(t, "success_template.html", "Hello, {{ .Var.Name }}! {{ .Config.Static }}/main.js")
	testutils.WriteFile(t, "failure_template.html", "Hello, {{ .Var.Name!")

	// Successful rendering scenario
	tpl := &gsk.Tpl{
		TemplatePath: "success_template.html",
		Variables:    map[string]string{"Name": "World"},
	}
	data, err := tpl.Render(gsk.DEFAULT_TEMPLATE_VARIABLES)
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World! /static/main.js", string(data))

	// Parsing error scenario
	tpl = &gsk.Tpl{
		TemplatePath: "failure_template.html",
	}
	_, err = tpl.Render(gsk.DEFAULT_TEMPLATE_VARIABLES)
	assert.Error(t, err)

}
