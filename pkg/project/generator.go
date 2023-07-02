package project

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func Generate(config *Config) (err error) {

	config.DirTree = dirTree
	config.DirNames = dirNames

	// use existing repo name as package name
	// or initialize git repo
	err = initializePackageWithGit(config)
	if err != nil {
		log.Fatal("error initializing go package with git: ", err)
		return err
	}

	// run go mod init
	err = exec.Command("go", "mod", "init", config.PkgName).Run()
	if err != nil {
		log.Fatal("error initializing go module: ", err)
		return err
	}

	// create dirs
	CreateProjectStructure(dirList)

	err = CreateProjectFiles(config)
	if err != nil {
		log.Fatal("error creating project files: ", err)
		return err
	}

	log.Println("Running go mod tidy...")
	// run go mod tidy
	err = exec.Command("go", "mod", "tidy").Run()
	if err != nil {
		log.Fatal("error initializing go module: ", err)
		return err
	}

	return nil
}

func RandomName() string {
	adjectives := []string{"dusty", "shiny", "noisy", "happy", "chubby", "fluffy"}
	nouns := []string{"donuts", "bears", "kittens", "ducks", "apples", "oranges"}

	randSrc := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSrc)

	adjective := adjectives[randGen.Intn(len(adjectives))]
	noun := nouns[randGen.Intn(len(nouns))]

	return fmt.Sprintf("%s%s", adjective, noun)
}

func CreateProjectFiles(config *Config) (err error) {

	filesToCreate := [][]FileTemplate{
		cmdFileTemplate(config),
		entityFileTemplate(config),
		serviceFileTemplate(config),
		handlerFileTemplate(config),
		serverFileTemplate(config),
		storageFileTemplate(config),
	}

	for _, files := range filesToCreate {
		err = CreateFilesFromTemplate(files, config)
		if err != nil {
			return err
		}
	}

	return nil
}

func CreateFilesFromTemplate(files []FileTemplate, config *Config) (err error) {
	for _, file := range files {
		f, err := os.Create(file.Path)
		if err != nil {
			return err
		}

		err = file.Template.Execute(f, config)
		if err != nil {
			return err
		}
	}
	return nil
}
