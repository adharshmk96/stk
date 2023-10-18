package progen

import (
	"log"
	"os"
	"os/exec"
)

func IsGoModule() bool {
	_, err := os.Stat("go.mod")
	return !os.IsNotExist(err)
}

func runCommand(name string, args ...string) error {
	err := exec.Command(name, args...).Run()
	if err != nil {
		log.Fatalf("error running %s: %v", name, err)
	}
	return err
}

// openDirectory creates (if not exists) and changes the working directory to the given path.
func OpenDirectory(workDir string) error {
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return err
	}
	return os.Chdir(workDir)
}
