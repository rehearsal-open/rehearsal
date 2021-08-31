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

	"github.com/rehearsal-open/rehearsal/rehearsal-cli/cli"
	"github.com/rehearsal-open/rehearsal/task"
)

type (
	Frontend struct {
		logger *cli.Task
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
