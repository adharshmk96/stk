package commands

import (
	"bytes"
	"errors"
	"os/exec"
	"strings"
)

type GitCmd interface {
	RunCmd(args ...string) (string, error)
	Init() error
	AddRemote(remoteName, remoteUrl string) error
	GetRemoteOrigin() (string, error)
	Revparse(ref string) (string, error)

	IsRepo() bool
}

type gitCommands struct{}

func NewGitCmd() GitCmd {
	return &gitCommands{}
}

func (g *gitCommands) RunCmd(args ...string) (string, error) {
	stdout := bytes.Buffer{}
	stderr := bytes.Buffer{}

	c := exec.Command("git", args...)
	c.Stdout = &stdout
	c.Stderr = &stderr

	err := c.Run()
	if err != nil {
		return "", errors.New(stderr.String())
	}
	return stdout.String(), nil
}

// IsRepo returns true if current folder is a git repository.
func (g *gitCommands) IsRepo() bool {
	out, err := g.Revparse("--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

func (g *gitCommands) Init() error {
	_, err := g.RunCmd("init", ".")
	if err != nil {
		return err
	}

	return nil
}

func (g *gitCommands) AddRemote(remoteName, remoteUrl string) error {
	_, err := g.RunCmd("remote", "add", remoteName, remoteUrl)
	if err != nil {
		return err
	}

	return nil
}

func (g *gitCommands) GetRemoteOrigin() (string, error) {
	out, err := g.RunCmd("config", "--get", "remote.origin.url")
	if err != nil {
		return "", err
	}

	return out, nil
}

func (g *gitCommands) Revparse(ref string) (string, error) {
	out, err := g.RunCmd("rev-parse", ref)
	if err != nil {
		return "", err
	}

	return out, nil
}
