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
	"go.bug.st/serial"
)

func (t *__task) IsSupporting(elem task_element.Enum) bool {
	return [task_element.Len]bool{
		true, true, false,
	}[elem]
}

func (t *__task) ExecuteMain(args based.MainFuncArguments) error {

	var stdOut io.Writer

	if stdout, err := args.Writer(task_element.StdOut); err != nil {
		return errors.WithStack(err)
	} else {
		stdOut = stdout
	}

	listen.Listen(t, task_element.StdIn, t.Port, func(e error) {
		if err, ok := e.(serial.PortError); ok {
			switch err.Code() {
			case serial.PortClosed:
				t.StopMain()
			}
		}
	}, nil)
	closer := listen.SyncIoPipe(t.Port, stdOut, time.Millisecond, func(e error) {
		if err, ok := e.(serial.PortError); ok {
			switch err.Code() {
			case serial.PortClosed:
				t.StopMain()
			}
		}
	})

	go func() {

		<-t.close
		closer <- nil
		exitErr := t.Port.Close()

		args.Close(exitErr)
	}()

	return nil
}

func (t *__task) StopMain() {
	close(t.close)
}
