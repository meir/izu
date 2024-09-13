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

//go:embed template
var template string

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

	keys := map[string]string{}
	for _, match := range regex.FindAllStringSubmatch(string(content), -1) {
		keys[strings.ToLower(match[1])] = fmt.Sprintf("\"%s\": \"%s\"", strings.ToLower(match[1]), match[1])
	}

	key_array := []string{}
	for _, key := range keys {
		key_array = append(key_array, key)
	}
	slices.Sort(key_array)

	template = fmt.Sprintf(template, source, strings.Join(key_array, ", "))

	err = os.WriteFile("generated_keys.go", []byte(template), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
