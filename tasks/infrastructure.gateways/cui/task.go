// tasks/infrastructure.gateways/cui/task.go
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
	"os/exec"

	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/element"
	"github.com/rehearsal-open/rehearsal/tasks/basis.gateways"
	"github.com/rehearsal-open/rehearsal/tasks/buffer"
	"github.com/rehearsal-open/rehearsal/tasks/gateways"
	"github.com/rehearsal-open/rehearsal/tasks/infrastructure.gateways/cui.entity"
)

func Make(entity *entities.Task) (cuiTask gateways.Task, err error) {

	var (
		result task
	)

	if detail, ok := entity.Detail().(*cui.Detail); !ok {
		return nil, gateways.ErrInvalidTaskDetail
	} else {
		result = task{
			entity: entity,
			Detail: detail,
			stdin:  nil,
		}
	}

	result.Task = basis.MakeBasis(entity, &result)
	result.Cmd = exec.Command(result.Detail.Path, result.Detail.Args...)
	result.Cmd.Dir = result.Detail.Dir

	return &result, nil
}

func (cui *task) IsSupporting(elem element.TaskElement) bool {
	return [element.NumTaskElement]bool{
		true, true, true,
	}[elem]
}

func (cui *task) ExecuteMain(args basis.MainFuncArguments) error {

	var stdIn io.Writer

	// listener
	if stdin, err := cui.Cmd.StdinPipe(); err != nil {
		return err
	} else {
		stdIn = stdin
	}

	// writer
	if out, err := args.Writer(element.StdOut); err != nil {
		return err
	} else {
		cui.Cmd.Stdout = out
	}
	if out, err := args.Writer(element.StdErr); err != nil {
		return err
	} else {
		cui.Cmd.Stderr = out
	}

	callback := [element.NumTaskElement]basis.RecieveCallback{nil}
	callback[element.StdIn] = func(recieve buffer.Packet) error {
		if stdIn != nil {
			buffer := bytes.NewBuffer(make([]byte, 0))
			io.Copy(buffer, &recieve)
			if _, err := io.Copy(stdIn, buffer); err != nil {
				panic(err.Error())
			}
			return nil
		} else {
			panic("stdin finalized")
		}
	}

	if err := cui.Start(); err != nil {
		return err
	} else {

		args.ListenStart(callback)

		// start running element
		go func() {
			err := error(nil)
			err = cui.Cmd.Wait()
			if err != nil {
				panic(err.Error())
			}
			args.Close(err)
			cui.stdin = nil
		}()

		return nil
	}

}

func (cui *task) StopMain() {
	cui.Cmd.Process.Kill()
}
