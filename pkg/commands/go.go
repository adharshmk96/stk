package commands

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
)

type GoCmd interface {
	RunCmd(args ...string) (string, error)
	IsMod() bool

	ModInit(string) error
	ModTidy() error

	ModPackageName() (string, error)
}

type goCommands struct{}

func NewGoCmd() GoCmd {
	return &goCommands{}
}

func (g *goCommands) RunCmd(args ...string) (string, error) {
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	c := exec.Command("go", args...)
	c.Stdout = &stdout
	c.Stderr = &stderr

	err := c.Run()
	if err != nil {
		return "", errors.New(stderr.String())
	}
	return stdout.String(), nil
}

func (g *goCommands) IsMod() bool {
	_, err := os.Stat("go.mod")
	return !os.IsNotExist(err)
}

func (g *goCommands) ModInit(pkg string) error {
	_, err := g.RunCmd("mod", "init", pkg)
	if err != nil {
		return err
	}

	return nil
}

func (g *goCommands) ModTidy() error {
	_, err := g.RunCmd("mod", "tidy")
	if err != nil {
		return err
	}

	return nil
}

func (g *goCommands) ModPackageName() (string, error) {
	if !g.IsMod() {
		return "", errors.New("not a Go module")
	}

	out, err := Clean(g.RunCmd("go", "list", "-m"))
	if err != nil {
		return "", err
	}

	return out, nil
}
