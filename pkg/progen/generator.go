package progen

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/adharshmk96/stk/pkg/progen/tpl"
)

const MODULE_PLACEHOLDER = "ping"

type Generator struct {
	Config *Config
}

func NewGenerator(config *Config) *Generator {
	return &Generator{
		Config: config,
	}
}

func (g *Generator) GenerateProject() error {
	if !g.Config.IsGitRepo {
		log.Println("Initializing git repository...")
		if err := initGit(); err != nil {
			return err
		}
	}

	if !g.Config.IsGoModule {
		log.Println("Initializing go module...")
		if err := initGoMod(); err != nil {
			return err
		}
	}

	log.Println("Generating project files...")
	g.generateBoilerplate(tpl.ProjectTemplates)

	log.Println("Running go mod tidy...")
	if err := goModTidy(); err != nil {
		return err
	}

	return nil
}

func (g *Generator) GenerateModule() error {
	g.generateBoilerplate(tpl.ModuleTemplates)
	return nil
}

func (g *Generator) DeleteModule() error {
	return g.deleteModuleBoilerplate(tpl.ModuleTemplates)
}

func (g *Generator) generateBoilerplate(templates []tpl.Template) {
	for _, tf := range templates {
		tf.FilePath = formatModuleFilePath(tf.FilePath, g.Config)

		if err := createDirectoryForFile(tf.FilePath); err != nil {
			continue
		}

		if err := createAndExecuteTemplate(tf, g.Config); err != nil {
			continue
		}
	}
}

func (g *Generator) deleteModuleBoilerplate(templates []tpl.Template) error {
	for _, tf := range templates {
		tf.FilePath = formatModuleFilePath(tf.FilePath, g.Config)
		if err := os.Remove(tf.FilePath); err != nil {
			log.Fatalf("Failed to delete file %s: %v\n", tf.FilePath, err)
		}
	}
	return nil
}

func formatModuleFilePath(pathTemplate string, config *Config) string {
	return strings.ReplaceAll(pathTemplate, MODULE_PLACEHOLDER, config.ModName)
}

func createDirectoryForFile(filePath string) error {
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Fatalf("Failed to create directory for file %s: %v\n", filePath, err)
	}
	return err
}

func createAndExecuteTemplate(tf tpl.Template, config *Config) error {
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

func RandomName() string {
	nouns := []string{"apple", "ball", "cat", "dog", "elephant", "fish", "gorilla", "horse", "iguana", "jellyfish", "kangaroo"}
	adjectives := []string{"angry", "big", "cold", "dark", "fast", "good", "happy", "jolly", "kind", "little", "merry", "nice"}

	randSrc := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSrc)

	return fmt.Sprintf("%s%s", adjectives[randGen.Intn(len(adjectives))], nouns[randGen.Intn(len(nouns))])
}
