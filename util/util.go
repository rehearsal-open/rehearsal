package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/pkg/errors"
)

// Return whether file is exist
func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// Get object from Mark up language file, supported json(*.json) and yaml(*.yaml, *.yml).
func MarkUpFileLoad(filename string, result interface{}) (string, error) {
	buf, err := ioutil.ReadFile(filename)

	if err != nil {
		return "", errors.WithStack(err)
	}

	ext := filepath.Ext(filename)

	switch ext {
	case ".yaml":
	case ".yml":
		err = yaml.Unmarshal(buf, result)
		break
	case ".json":
		err = json.Unmarshal(buf, result)
		break
	default:
		return ext, errors.New("Invalid File Extension: " + ext)
	}

	return ext, nil
}
