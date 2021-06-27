package parser

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entity"
)

func LoadArgs(config *entity.Conf) error {
	args := os.Args

	switch args[1] {
	case "run":
		if path, err := filepath.Abs(args[2]); err != nil {
			return errors.WithMessage(err, "invalid file: "+path)
		} else if _, err := os.Stat(path); err != nil {
			return errors.WithMessage(err, "cannot found file: "+path)
		} else {
			switch filepath.Ext(path) {
			case ".yml":
				return errors.WithMessage(Yaml(path, config), "yaml parse failed.")
			case ".yaml":
				return errors.WithMessage(Yaml(path, config), "yaml parse failed.")
			default:
				return errors.New("not supported file extension: " + path)
			}
		}
	default:
		return errors.New("invalid command: " + args[1])
	}
}
