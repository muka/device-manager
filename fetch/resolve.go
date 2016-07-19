package fetch

import (
	"os/user"
	"path/filepath"
	"strings"
)

// ResolvePath transform a path to an absolute path
func ResolvePath(path string) (string, error) {

	if filepath.IsAbs(path) {
		return path, nil
	}

	usr, err := user.Current()
	if err != nil {
		return path, err
	}

	dir := usr.HomeDir + "/"

	if path[:2] == "~/" {
		path = strings.Replace(path, "~/", dir, 1)
	}

	return filepath.Abs(path)
}
