// rehearsal-cli/run.go
// Copyright (C) 2021 Kasai Koji

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

package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/engine"
	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/parser/yaml"
	"github.com/rehearsal-open/rehearsal/rehearsal-cli/cli"
	"github.com/rehearsal-open/rehearsal/task/impl/cui"
	"github.com/rehearsal-open/rehearsal/task/impl/serial"
	"github.com/rehearsal-open/rehearsal/task/maker"
)

var (
	SupportedTasks *maker.Maker
	Wd             string
)

func init_run() {
	SupportedTasks = &maker.Maker{}
	SupportedTasks.RegisterMaker("cui", &cui.MakeCollection)
	SupportedTasks.RegisterMaker("serial", &serial.MakeCollection)
	Wd, _ = os.Getwd()
}

func Run(confFile string) error {

	if abs, err := filepath.Abs(confFile); err != nil {
		return errors.WithMessage(err, "config file path is invalid")
	} else {
		confFile = abs
	}

	parser := yaml.Parser{
		Path:        confFile,
		DetailMaker: SupportedTasks,
	}

	entity, en := entities.Rehearsal{
		DefaultDir: filepath.Dir(parser.Path),
	}, engine.Rehearsal{}

	if err := parser.Parse(&entity); err != nil {
		return errors.WithStack(err)
	} else if logger, err := cli.MakeTask(&entity); err != nil {
		return errors.WithStack(err)
	} else {
		frontend := Frontend{
			logger: logger,
		}

		SupportedTasks.Frontend = &frontend

		if err := en.Reset(&entity, SupportedTasks, &frontend); err != nil {
			return errors.WithStack(err)
		}
	}

	if err := en.Execute(); err != nil {
		return errors.WithStack(err)
	}

	time.Sleep(time.Second)

	return nil
}
