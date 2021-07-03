// load/load.go
// Copyright (C) 2021  Kasai Koji

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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
	case "about":
		conf.Command = "about"
		return &conf, nil
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

		conf.Dir = filepath.Join(filepath.Dir(conf.YamlPath), conf.Dir)
	}

	return nil
}
