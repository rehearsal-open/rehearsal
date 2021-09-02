// rehearsal-cli/objects.go
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
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/parser"
	"github.com/rehearsal-open/rehearsal/rehearsal-cli/cli"
	"github.com/rehearsal-open/rehearsal/task"
	"github.com/rehearsal-open/rehearsal/task/based"
	"github.com/rehearsal-open/rehearsal/task/impl/cui"
	"github.com/streamwest-1629/convertobject"
)

type (
	Frontend struct {
		logger *cli.Task
		config *Config
	}

	Config struct {
		BaseDir string `map-to:"dir"`
	}
)

func (f *Frontend) LoggerTask() task.Task {
	return f.logger
}

func (f *Frontend) Log(flag int, msg string) {
	fmt.Println("[INFO]: ", msg)
}

func (f *Frontend) Select(msg string, options []string) int {
	fmt.Println("[SELECT]: " + msg + "(type number).")
	for i, op := range options {
		fmt.Println("  [", i+1, "]: ", op)
	}
	selected := 0
	for {
		fmt.Scanf("%d", &selected)
		if selected > 0 && selected <= len(options) {
			f.Log(0, fmt.Sprint("OK, you selected '", options[selected-1], "'."))
			return selected - 1
		} else {
			f.Log(0, "you selected invalid number, try again.")
		}
	}
}

// Set default value
func (f *Frontend) InitConfig(src parser.MappingType) error {

	parseRes := Config{}

	if err := convertobject.DirectConvert(src, &parseRes); err != nil {
		return errors.WithStack(err)
	}

	// initialize filepath
	if parseRes.BaseDir != "" {
		if filepath.IsAbs(parseRes.BaseDir) {
			f.config.BaseDir = parseRes.BaseDir
		} else {
			f.config.BaseDir = filepath.Join(parseRes.BaseDir)
		}
	}

	// set default execute directory
	cui.DefaultDir = f.config.BaseDir

	f.logger = &cli.Task{}
	f.logger.Task = based.MakeBasis(&entities.Task{}, f.logger)

	return nil

}
