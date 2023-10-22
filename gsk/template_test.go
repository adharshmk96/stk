package gsk_test

import (
	"os"
	"testing"

	"github.com/adharshmk96/stk/gsk"
	"github.com/adharshmk96/stk/testutils"
	"github.com/stretchr/testify/assert"
)

var templateContent = `<html>
	<head>
		<title>{{ .Title }}</title>
	</head>
	<body>
		<h1>{{ .Title }}</h1>
		<p>{{ .Body }}</p>
	</body>
</html>`

func TestTemplateResponseRender(t *testing.T) {
	t.Run("renders template with data", func(t *testing.T) {
		tempFile, removeDir := testutils.CreateTempFile(t, templateContent)
		defer removeDir()

		data := struct {
			Title string
			Body  string
		}{
			Title: "Hello",
			Body:  "World",
		}

		templateResp := gsk.Tpl{
			TemplatePath: tempFile,
			Variables:    data,
		}

		rendered, err := templateResp.Render()
		assert.NoError(t, err)

		expected := `<html>
	<head>
		<title>Hello</title>
	</head>
	<body>
		<h1>Hello</h1>
		<p>World</p>
	</body>
</html>`

		assert.Equal(t, expected, string(rendered))
	})

	t.Run("execute error is returned", func(t *testing.T) {
		tempFile, _ := testutils.CreateTempFile(t, templateContent)

		data := struct {
			Title string
			Body  string
		}{
			Title: "Hello",
			Body:  "World",
		}

		templateResp := gsk.Tpl{
			TemplatePath: tempFile,
			Variables:    data,
		}

		// Remove template file to force execute error
		err := os.Remove(tempFile)
		assert.NoError(t, err)

		_, err = templateResp.Render()
		assert.Error(t, err)
	})
}
