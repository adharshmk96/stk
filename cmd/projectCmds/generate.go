/*
Copyright © 2023 Adharsh M dev@adharsh.in
*/
package projectCmds

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/adharshmk96/stk/pkg/project"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getWorkDirFromArg(args []string) string {
	if len(args) == 0 {
		return "."
	}

	return args[0]
}

func getRepoName() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	repoUrl := string(out)
	repoUrl = strings.TrimSuffix(repoUrl, ".git\n")
	repoUrl = strings.ReplaceAll(repoUrl, "https://", "")
	repoUrl = strings.ReplaceAll(repoUrl, "git@", "")
	repoUrl = strings.ReplaceAll(repoUrl, ":", "/")

	return repoUrl, nil
}

func getAppNameFromPkgName(s string) string {
	lastSlash := strings.LastIndex(s, "/")
	lastPart := s
	if lastSlash != -1 {
		lastPart = s[lastSlash+1:]
	}
	return strings.ReplaceAll(lastPart, "-", "")
}

func openDirectory(workDir string) error {
	os.MkdirAll(workDir, 0755)

	err := os.Chdir(workDir)
	if err != nil {
		return err
	}
	return nil
}

// init with git after cd
// git init
// git remote add origin <your_repository_url>
// git fetch
// git checkout -t origin/<your_branch_name>

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a project with gsk and clean architecture.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Generating project files...")

		workdir := getWorkDirFromArg(args)
		pkg := viper.GetString("project.package")
		app := viper.GetString("project.app")

		err := openDirectory(workdir)
		if err != nil {
			log.Fatal(err)
			return
		}

		packageNameFromGit, err := getRepoName()
		if err == nil && packageNameFromGit != "" {
			// if a git repo exists, use the repo name as package name
			log.Println("using existing git repo name: ", packageNameFromGit)
			pkg = packageNameFromGit
			app = getAppNameFromPkgName(pkg)
		} else {
			randomName := project.RandomName()
			if pkg == "" {
				pkg = randomName
			}

			if app == "" {
				app = randomName
			}
		}

		app = strings.ReplaceAll(app, "-", "")
		app = strcase.ToLowerCamel(app)

		config := &project.Config{
			RootPath: workdir,
			PkgName:  pkg,
			AppName:  app,
		}

		err = project.Generate(config)
		if err != nil {
			log.Fatal(err)
			return
		}

		log.Println("Project generated successfully.")
	},
}
