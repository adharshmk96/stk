package project

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/adharshmk96/stk/pkg/project/tpl"
)

const MODULE_PLACEHOLDER = "ping"

func GenerateProjectBoilerplate(ctx *Context) error {
	if !ctx.IsGitRepo {
		fmt.Println("initializing git repository...")
		err := ctx.GitCmd.Init()
		if err != nil {
			return err
		}
	}

	if !ctx.IsGoModule {
		fmt.Println("initializing go module...")
		err := ctx.GoCmd.ModInit(ctx.PackageName)
		if err != nil {
			return err
		}
	}

	fmt.Println("generating project files...")
	templateConfig := GetTemplateConfig(ctx, "")
	err := generateBoilerplate(tpl.ProjectTemplates, templateConfig)
	if err != nil {
		return err
	}
	for _, module := range ctx.Modules {
		templateConfig := GetTemplateConfig(ctx, module)
		err = generateBoilerplate(tpl.ModuleTemplates, templateConfig)
		if err != nil {
			return err
		}
	}

	fmt.Println("running go mod tidy...")
	err = ctx.GoCmd.ModTidy()
	if err != nil {
		return err
	}

	return nil
}

func GenerateModuleBoilerplate(ctx *Context, module string) error {
	fmt.Println("generating module files...")
	templateConfig := GetTemplateConfig(ctx, module)
	return generateBoilerplate(tpl.ModuleTemplates, templateConfig)
}

func DeleteModuleBoilerplate(ctx *Context, module string) error {
	fmt.Println("deleting module files...")
	templateConfig := GetTemplateConfig(ctx, module)
	return deleteBoilerplate(tpl.ModuleTemplates, templateConfig)
}

func generateBoilerplate(templates []tpl.Template, config *TemplateConfig) error {
	for _, tf := range templates {
		tf.FilePath = formatModuleFilePath(tf.FilePath, config)

		if err := createDirectoryForFile(tf.FilePath); err != nil {
			return err
		}

		if err := createAndExecuteTemplate(tf, config); err != nil {
			return err
		}
	}
	return nil
}

func deleteBoilerplate(templates []tpl.Template, config *TemplateConfig) error {
	for _, tf := range templates {
		tf.FilePath = formatModuleFilePath(tf.FilePath, config)

		err := os.Remove(tf.FilePath)
		if err != nil {
			return err
		}
	}
	return nil
}

func formatModuleFilePath(pathTemplate string, config *TemplateConfig) string {
	return strings.ReplaceAll(pathTemplate, MODULE_PLACEHOLDER, config.ModName)
}

func createDirectoryForFile(path string) error {
	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatalf("failed to create directory for file %s: %v\n", path, err)
	}
	return err
}

func createAndExecuteTemplate(tf tpl.Template, config *TemplateConfig) error {
	f, err := os.Create(tf.FilePath)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v\n", tf.FilePath, err)
		return err
	}
	defer f.Close()

	tpl := template.Must(template.New(tf.FilePath).Parse(tf.Content))
	err = tpl.Execute(f, config)
	if err != nil {
		log.Fatalf("Failed to execute template for file %s: %v\n", tf.FilePath, err)
	}
	return err
}
