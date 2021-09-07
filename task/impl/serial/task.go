// task/impl/serial/task.go
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

package serial

import (
	"io"
	"time"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/task/based"
	"github.com/rehearsal-open/rehearsal/task/wrapper/listen"
)

func (serial *__task) IsSupporting(elem task_element.Enum) bool {
	return [task_element.Len]bool{
		true, true, false,
	}[elem]
}

func (serial *__task) ExecuteMain(args based.MainFuncArguments) error {

	var stdOut io.Writer

	if stdout, err := args.Writer(task_element.StdOut); err != nil {
		return errors.WithStack(err)
	} else {
		stdOut = stdout
	}

	listen.Listen(serial, task_element.StdIn, serial.Port, nil, nil)
	closer := listen.SyncIoPipe(serial.Port, stdOut, time.Millisecond, nil)

	go func() {

		exitErr := serial.Port.Close()
		closer <- nil

		args.Close(exitErr)
	}()

	return nil
}

func (serial *__task) StopMain() {
	close(serial.close)
}
