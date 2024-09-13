package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"
)

const source = "https://raw.githubusercontent.com/xkbcommon/libxkbcommon/master/include/xkbcommon/xkbcommon-keysyms.h"

//go:embed template_map
var template_map string

//go:embed template_switch
var template_switch string

//go:embed template_array
var template_array string

var regex = regexp.MustCompile(`\#define XKB_KEY_([a-zA-Z_0-9]+)\s`)

func main() {
	response, err := http.Get(source)
	if err != nil {
		log.Fatal(err)
	}

	content, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	keys := []string{}
	for _, match := range regex.FindAllStringSubmatch(string(content), -1) {
		if match[1] == "" {
			continue
		}
		keys = append(keys, fmt.Sprintf("\"%s\"", match[1]))
	}

	slices.Sort(keys)
	keys = slices.Compact(keys)

	generate_switch(keys)
	generate_array(keys)
	generate_map(keys)
}

func generate_map(keys []string) {
	for i, key := range keys {
		keys[i] = key + ": true"
	}

	template_map = fmt.Sprintf(template_map, source, strings.Join(keys, ", "))

	err := os.WriteFile("generated_keys_map.go", []byte(template_map), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func generate_switch(keys []string) {
	template_switch = fmt.Sprintf(template_switch, source, strings.Join(keys, ", "))

	err := os.WriteFile("generated_keys_switch.go", []byte(template_switch), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func generate_array(keys []string) {
	template_array = fmt.Sprintf(template_array, source, strings.Join(keys, ", "))
	err := os.WriteFile("generated_keys_array.go", []byte(template_array), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
