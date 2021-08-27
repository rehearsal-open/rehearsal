// rehearsal-cli/cli/objects.go
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

package cli

import (
	"fmt"

	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/task"
)

const (
	ForeRed     = "\x1b[31m"
	ForeGreen   = "\x1b[32m"
	ForeYellow  = "\x1b[33m"
	ForeBlue    = "\x1b[34m"
	ForeMagenta = "\x1b[35m"
	ForeSyan    = "\x1b[36m"
	ForeReset   = "\x1b[39m"
	BackRed     = "\x1b[41m"
	BackGreen   = "\x1b[42m"
	BackYellow  = "\x1b[43m"
	BackBlue    = "\x1b[44m"
	BackMagenta = "\x1b[45m"
	BackSyan    = "\x1b[46m"
	BackReset   = "\x1b[49m"
)

type (
	Cli struct {
		Entity *entities.Rehearsal
	}
)

func (c *Cli) LoggerTask() task.Task {
	task, _ := MakeTask(c.Entity)
	return task
}

func (c *Cli) Log(flag int, msg string) {
	fmt.Println(msg)
}
