package project

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/adharshmk96/stk/pkg/project/tpl"
)

func GenerateProject(config *ProjectConfig) error {
	// use existing repo name as package name
	// or initialize git repo
	log.Println("Initializing git repository...")
	err := initializePackageWithGit(config)
	if err != nil {
		log.Fatal("error initializing go package with git: ", err)
		return err
	}

	// run go mod init
	log.Println("Running go mod init...")
	err = exec.Command("go", "mod", "init", config.PkgName).Run()
	if err != nil {
		log.Fatal("error initializing go module: ", err)
		return err
	}

	// create boilerplate
	log.Println("Generating boilerplate...")
	generateBoilerplate(config)

	// run go mod tidy
	log.Println("Running go mod tidy...")
	err = exec.Command("go", "mod", "tidy").Run()
	if err != nil {
		log.Fatal("error running go mod tidy: ", err)
		return err
	}

	return nil
}

func formatModuleFilePath(pathTemplate string, modConfig *ModuleConfig) string {
	filePath := strings.ReplaceAll(pathTemplate, "{{ .ModName }}", modConfig.ModName)
	return filePath
}

func GenerateModule(modConfig *ModuleConfig) error {
	log.Println("Adding boilerplate for module...")
	templates := tpl.ModuleTemplates
	for _, tf := range templates {
		tf.FilePath = formatModuleFilePath(tf.FilePath, modConfig)
		dir := filepath.Dir(tf.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory for file %s: %v\n", tf.FilePath, err)
			return err
		}

		f, err := os.Create(tf.FilePath)
		if err != nil {
			log.Fatalf("Failed to create file %s: %v\n", tf.FilePath, err)
			return err
		}
		defer f.Close()

		tpl := template.Must(template.New(tf.FilePath).Parse(tf.Content))

		if err := tpl.Execute(f, modConfig); err != nil {
			log.Fatalf("Failed to execute template for file %s: %v\n", tf.FilePath, err)
			return err
		}

	}

	return nil
}

func generateBoilerplate(config *ProjectConfig) {
	templates := tpl.SingleModTemplates

	for _, tf := range templates {
		dir := filepath.Dir(tf.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Fatalf("Failed to create directory for file %s: %v\n", tf.FilePath, err)
			continue
		}

		f, err := os.Create(tf.FilePath)
		if err != nil {
			log.Fatalf("Failed to create file %s: %v\n", tf.FilePath, err)
			continue
		}
		defer f.Close()

		tpl := template.Must(template.New(tf.FilePath).Parse(tf.Content))

		if err := tpl.Execute(f, config); err != nil {
			log.Fatalf("Failed to execute template for file %s: %v\n", tf.FilePath, err)
			continue
		}

	}
}

func RandomName() string {
	nouns := []string{
		"apple",
		"ball",
		"cat",
		"dog",
		"elephant",
		"fish",
		"gorilla",
		"horse",
		"iguana",
		"jellyfish",
		"kangaroo",
	}

	adjectives := []string{
		"angry",
		"big",
		"cold",
		"dark",
		"fast",
		"good",
		"happy",
		"jolly",
		"kind",
		"little",
		"merry",
		"nice",
	}

	randSrc := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSrc)

	adjective := adjectives[randGen.Intn(len(adjectives))]
	noun := nouns[randGen.Intn(len(nouns))]

	return fmt.Sprintf("%s%s", adjective, noun)

}
