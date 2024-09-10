package util

import (
	"os"
	"regexp"
	"strings"
)

const LogsDir = "var/log"

// Path is a function which builts any path from root dir.
func Path(additionalPath string) (string, error) {
	root, err := os.Getwd()
	if err != nil {
		return "", err
	}

	b := strings.Builder{}
	b.WriteString(root)
	b.WriteRune(os.PathSeparator)
	b.WriteString(additionalPath)

	return replaceDuplicatedSlashes(b.String()), nil
}

func LogsPath(additionalPath string) (string, error) {
	logsDir, err := Path(LogsDir)
	if err != nil {
		return "", err
	}

	return replaceDuplicatedSlashes(logsDir + string(os.PathSeparator) + additionalPath), nil
}

func replaceDuplicatedSlashes(path string) string {
	return regexp.MustCompile(`/+`).ReplaceAllString(path, string(os.PathSeparator))
}
