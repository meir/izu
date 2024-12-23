package izu

import (
	"embed"
	"fmt"
	"os"
	"strings"
)

// Formatter is the interface that should be implemented for all hotkey formatters
// These are used to interface with a specific formatter language
type Formatter interface {
	ParseString([]byte) ([]byte, error)
	ParseFile(string) ([]byte, error)
}

//go:embed formatters/*
var formatters embed.FS

// GetVersion will return the version stated in the VERSION file
func GetVersion() string {
	version, err := formatters.ReadFile("formatters/VERSION")
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(version))
}

func GetFormatterFile(language, system string) ([]byte, error) {
	content, err := formatters.ReadFile(fmt.Sprintf("formatters/%s/%s.lua", language, system))
	if err != nil {
		// system might be a file path instead of a system name
		content, err = os.ReadFile(system)
	}
	return content, err
}
