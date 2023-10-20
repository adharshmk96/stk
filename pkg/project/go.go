package project

import (
	"errors"
	"os"
)

func InitGoMod(packageName string) error {
	_, err := RunCmd("go", "mod", "init", packageName)
	if err != nil {
		return err
	}

	return nil
}

func GoModTidy() error {
	_, err := RunCmd("go", "mod", "tidy")
	if err != nil {
		return err
	}

	return nil
}

func IsGoModule() bool {
	_, err := os.Stat("go.mod")
	return !os.IsNotExist(err)
}

func GoModPackageName() (string, error) {
	if !IsGoModule() {
		return "", errors.New("not a Go module")
	}

	out, err := Clean(RunCmd("go", "list", "-m"))
	if err != nil {
		return "", err
	}

	return out, nil
}
