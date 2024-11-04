package izu

import "strings"

var keys = map[string]string{}

// CapitalizeKey will change the key to the capitalized version if it exists in the generated map with xkb keys
//
//go:generate go run ./gen/gen.go
func CapitalizeKey(binding [][]string) [][]string {
	output := [][]string{}

	for _, bind := range binding {
		line := []string{}
		for _, key := range bind {
			if value, ok := keys[strings.ToLower(key)]; ok {
				line = append(line, value)
			} else {
				line = append(line, key)
			}
		}
		output = append(output, line)
	}

	return output
}
