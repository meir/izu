package izu

import (
	"io"
	"os"
	"path"
	"strings"

	"github.com/meir/izu/pkg/izu"
)

func GetFormatter(name string) (data []byte, err error) {
	if !strings.HasSuffix(name, ".lua") {
		name += ".lua"
	}

	if file, err := izu.Formatters.Open(path.Join("formatters", name)); err == nil {
		// load file from embedded resources
		data, err = io.ReadAll(file)
	} else if os.IsNotExist(err) {
		// load file from disk
		if _, err = os.Stat(name); err != nil {
			return nil, os.ErrNotExist
		}

		data, err = os.ReadFile(name)
	}

	return
}
