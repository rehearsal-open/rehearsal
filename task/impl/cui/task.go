// task/impl/cui/task.go
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

package cui

import (
	"io"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/task"
	"github.com/rehearsal-open/rehearsal/task/based"
	"github.com/rehearsal-open/rehearsal/task/buffer"
	"github.com/rehearsal-open/rehearsal/task/impl/cui.entity"
)

func Make(entity *entities.Task) (t task.Task, err error) {

	result := __task{}

	if detail, ok := entity.Detail.(*cui.Detail); !ok {
		panic("invalid detail objects type")
	} else {
		result.Detail = detail
	}

	result.Task = based.MakeBasis(entity, &result)

	result.Cmd = exec.Command(result.Detail.Path, result.Detail.Args...)
	result.Cmd.Dir = result.Detail.Dir

	return &result, nil
}

func (cui *__task) IsSupporting(elem task_element.Enum) bool {
	return [task_element.Len]bool{
		true, true, true,
	}[elem]
}

func (cui *__task) ExecuteMain(args based.MainFuncArguments) error {

	var stdin io.Writer

	// listener
	if in, err := cui.Cmd.StdinPipe(); err != nil {
		return errors.WithStack(err)
	} else {
		stdin = in
	}

	// writer
	if out, err := args.Writer(task_element.StdOut); err != nil {
		return errors.WithStack(err)
	} else {
		cui.Cmd.Stdout = out
	}
	if out, err := args.Writer(task_element.StdErr); err != nil {
		return errors.WithStack(err)
	} else {
		cui.Cmd.Stderr = out
	}

	callback := [task_element.Len]based.Reciever{nil}
	callback[task_element.StdIn] = func(recieved *buffer.Packet) {
		if stdin != nil {
			defer recieved.Close()
			if _, err := io.Copy(stdin, recieved); err != nil {
				// todo: error manage
			}
		} else {
			panic("stdin finalized")
		}
	}

	if err := cui.Start(); err != nil {
		return errors.WithStack(err)
	}

	args.ListenStart(callback)

	// start running element
	go func() {
		args.Close(cui.Cmd.Wait())
		stdin = nil
	}()

	return nil
}

func (cui *__task) StopMain() {
	cui.Cmd.Process.Kill()
}
