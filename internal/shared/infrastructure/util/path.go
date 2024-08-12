package util

import (
	"os"
	"strings"
)

// Path is a function which builts any path from root dir.
func Path(additionalPath string) (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", err
	}

	b := strings.Builder{}
	b.WriteString(root)
	b.WriteString(additionalPath)
	fullPath := b.String()

	return fullPath, nil
}
