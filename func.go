package tools

import (
	"os"
	"os/exec"
	"path/filepath"
)

// AbsolutePath get execute binary path
func AbsolutePath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}

	return filepath.Dir(path) + "/", nil
}
