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
	"bytes"
	"io"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/task/based"
)

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

	callback := [task_element.Len]based.ImplCallback{nil}
	callback[task_element.StdIn] = based.MakeImplCallback(func(_ *entities.Element, b []byte) {
		if stdin != nil {
			if _, err := io.Copy(stdin, bytes.NewBuffer(b)); err != nil {
				panic(err.Error())
				// todo: error manage
			}
		}
	}, based.DefaultOnFinal)

	if err := cui.Start(); err != nil {
		return errors.WithStack(err)
	}

	args.ListenStart(callback)

	// start running element
	go func() {
		err := cui.Cmd.Wait()
		args.Close(err)
		stdin = nil
	}()

	return nil
}

func (cui *__task) StopMain() {
	cui.Cmd.Process.Kill()
}
