package envparser

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type ExtensionError struct {
	path string
	err  string
}

func (e *ExtensionError) Error() string { return fmt.Sprintf("path: %v %v", e.path, e.err) }

func IsCorrectExtension(path string) bool {
	const (
		DotEnvExtension    = ".env"
		DotEnvExtensionLen = 4
	)
	l := len(path)
	if l < DotEnvExtensionLen+1 {
		return false
	}
	for p, e := l-1, 3; p >= l-DotEnvExtensionLen; p, e = p-1, e-1 {
		if path[p] != DotEnvExtension[e] {
			return false
		}
	}
	return true
}

func Load(envFilePath string) error {
	if !IsCorrectExtension(envFilePath) {
		return &ExtensionError{path: envFilePath, err: "Expected .env file"}
	}
	f, err := os.Open(envFilePath)
	if err != nil {
		log.Fatalf("Load(). %s", err)
	}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		pair := strings.SplitN(scanner.Text(), "=", 2)
		if len(pair) != 2 {
			return &ExtensionError{path: envFilePath, err: "Expected lines with key=value pairs"}
		}
		err = os.Setenv(pair[0], pair[1])
		if err != nil {
			log.Fatal("Error: os.Setenv")
		}
	}
	return nil
}
