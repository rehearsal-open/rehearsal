// rehearsal-cli/flag.go
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
	"github.com/rehearsal-open/rehearsal/rehearsal-cli/cli"
	. "github.com/urfave/cli"
)

var (
	app *App
)

// initialize command
func init() {

	// initialize command and subcommand
	app = NewApp()
	app.Name = "rehearsal-cli"
	app.Description = "rehearsal-cli is testing-environment developping tool. More info: https://rehearsal-open.github.io"

	app.Commands = []Command{
		{
			// `rehearsal-cli run` command
			Name:        "run",
			Description: "rehearsal-cli executes tasks with template file.",
			Flags: []Flag{
				BoolFlag{ // plain text flag
					Name:  "plain",
					Usage: "rehearsal-cli prints log with plain text",
				},
				StringFlag{ // yaml file path
					Name:     "path, p",
					Usage:    "[REQUIRED] rehearsal-cli uses this file to execute and relation tasks",
					Required: true,
				},
			},
			Action: func(c *Context) {

				// execute rehearsal task
				cli.IsPlain = c.Bool("plain")
				if err := Run(c.String("path")); err != nil {
					println("ERROR: ", err.Error())
				}
			},
		},
		{
			// `rehearsal-cli version` command
			Name:        "version",
			Aliases:     []string{"-v"},
			Description: "print version",
			Action: func(c *Context) error {

				// print version
				Version()
				return nil
			},
		},
	}
}
