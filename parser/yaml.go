package parser

import (
	"os"

	"gopkg.in/yaml.v3"

	"github.com/pkg/errors"

	"github.com/rehearsal-open/rehearsal/entity"
)

func Yaml(filename string, config *entity.Conf) error {
	if fs, err := os.Open(filename); err != nil {
		return errors.WithStack(err)
	} else {
		defer fs.Close()

		decoder := yaml.NewDecoder(fs)
		if err := decoder.Decode(config); err != nil {
			return errors.WithStack(err)
		}

		return nil
	}
}
