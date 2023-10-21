package git

import (
	"strings"
)

// IsRepo returns true if current folder is a git repository.
func IsRepo() bool {
	out, err := Revparse("--is-inside-work-tree")
	return err == nil && strings.TrimSpace(out) == "true"
}

func Init() error {
	_, err := RunCmd("init", ".")
	if err != nil {
		return err
	}

	return nil
}

func AddRemote(remoteName, remoteUrl string) error {
	_, err := RunCmd("remote", "add", remoteName, remoteUrl)
	if err != nil {
		return err
	}

	return nil
}

func GetRemoteOrigin() (string, error) {
	out, err := RunCmd("config", "--get", "remote.origin.url")
	if err != nil {
		return "", err
	}

	return out, nil
}

func Revparse(ref string) (string, error) {
	out, err := RunCmd("rev-parse", ref)
	if err != nil {
		return "", err
	}

	return out, nil
}
