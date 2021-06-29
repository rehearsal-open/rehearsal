package load

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/rehearsal-open/rehearsal/entity"
)

func Load() (*entity.Config, error) {

	args := os.Args
	conf := entity.Config{}

	switch strings.ToLower(args[1]) {
	case "run":
		conf.Command = "run"
		if abs, err := filepath.Abs(args[2]); err != nil {
			return nil, errors.WithStack(err)
		} else {
			conf.YamlPath = abs
		}
		return &conf, errors.WithStack(loadConfigYaml(&conf))

	default:
		return nil, errors.New("unknown command: " + args[1])
	}
}

func loadConfigYaml(conf *entity.Config) error {

	if f, err := os.Open(conf.YamlPath); err != nil {
		return errors.WithStack(err)
	} else {
		defer f.Close()
		d := yaml.NewDecoder(f)

		if err := d.Decode(conf); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
