package gsk_test

import (
	"os"
	"testing"

	"github.com/adharshmk96/stk/gsk"
	"github.com/stretchr/testify/assert"
)

// Helper function to create temporary test template files.
func createTestFile(name string, content string) {
	f, _ := os.Create(name)
	defer f.Close()
	f.WriteString(content)
}

func TestMain(m *testing.M) {
	// Create temporary template files
	createTestFile("success_template.html", "Hello, {{.Name}}!")
	createTestFile("failure_template.html", "Hello, {{.Name!")
	// Run tests
	code := m.Run()
	// Cleanup
	os.Remove("success_template.html")
	os.Remove("failure_template.html")
	// Exit
	os.Exit(code)
}

func TestRender(t *testing.T) {
	// Successful rendering scenario
	tpl := &gsk.Tpl{
		TemplatePath: "success_template.html",
		Variables:    map[string]string{"Name": "World"},
	}
	data, err := tpl.Render()
	assert.NoError(t, err)
	assert.Equal(t, "Hello, World!", string(data))

	// Parsing error scenario
	tpl = &gsk.Tpl{
		TemplatePath: "failure_template.html",
	}
	_, err = tpl.Render()
	assert.Error(t, err)

}
