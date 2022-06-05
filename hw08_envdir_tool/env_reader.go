package main

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

func readFile(filename string) (EnvValue, error) {
	file, err := os.Open(filename)
	if err != nil {
		return EnvValue{}, err
	}
	defer file.Close()

	br := bufio.NewReader(file)
	line, err := br.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return EnvValue{}, err
	}

	line = strings.TrimRight(line, "\t\n ")
	line = strings.ReplaceAll(line, "\x00", "\n")

	needRemove := len(line) == 0

	return EnvValue{Value: line, NeedRemove: needRemove}, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	environments := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return environments, err
	}

	var filename string
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filename = file.Name()
		environments[filename], err = readFile(filepath.Join(dir, filename))
		if err != nil {
			return environments, err
		}
	}

	return environments, nil
}
