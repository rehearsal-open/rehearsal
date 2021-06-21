package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/pkg/errors"
)

func FileExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func MarkUpFileLoad(filename string) (string, map[string]interface{}, error) {
	buf, err := ioutil.ReadFile(filename)

	if err != nil {
		return "", nil, errors.WithStack(err)
	}

	var result map[string]interface{}
	ext := filepath.Ext(filename)

	switch ext {
	case ".yaml":
	case ".yml":
		err = yaml.Unmarshal(buf, &result)
		break
	case ".json":
		err = json.Unmarshal(buf, &result)
		break
	default:
		return ext, nil, errors.New("Invalid File Extension: " + ext)
	}

	return ext, result, nil
}
