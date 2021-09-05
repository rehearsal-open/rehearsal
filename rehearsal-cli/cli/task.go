// rehearsal-cli/cli/task.go
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
	"strings"
	"time"

	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/task/based"
)

type (
	Task struct {
		based.Task
		Format   map[string]string
		reciever chan string
		close    chan error
	}
)

func (task *Task) SetTaskEntities(entity *entities.Rehearsal) error {

	maxNameLen := 0
	task.Format = map[string]string{}

	colorSet := [...]string{
		ForeRed, ForeGreen, ForeYellow, ForeBlue, ForeMagenta, ForeSyan,
	}

	// get max length of logging task's names
	entity.ForeachTask(func(idx int, entity *entities.Task) error {

		if stdOut := entity.Element[task_element.StdOut]; stdOut.WriteLog {
			name := stdOut.Fullname()
			if stdOut.WriteLog {
				if nameSize := len(name); nameSize > maxNameLen {
					maxNameLen = nameSize
				}
			}
		}
		if stdErr := entity.Element[task_element.StdErr]; stdErr.WriteLog {
			name := stdErr.Fullname()
			if stdErr.WriteLog {
				if nameSize := len(name); nameSize > maxNameLen {
					maxNameLen = nameSize
				}
			}
		}
		return nil
	})

	idx := 0

	entity.ForeachTask(func(_ int, entity *entities.Task) error {

		if stdOut := entity.Element[task_element.StdOut]; stdOut.WriteLog {
			name := stdOut.Fullname()
			if stdOut.WriteLog {
				task.Format[name] = colorSet[idx%len(colorSet)] + BackReset + name + " (Std Out)" + strings.Repeat(" ", maxNameLen-len(name)) + " : "
				idx++
			}
		}
		if stdErr := entity.Element[task_element.StdErr]; stdErr.WriteLog {
			name := stdErr.Fullname()
			if stdErr.WriteLog {
				task.Format[name] = colorSet[idx%len(colorSet)] + BackReset + name + " (Std Err)" + strings.Repeat(" ", maxNameLen-len(name)) + " : "
				idx++
			}
		}

		return nil
	})

	return nil
}

func (t *Task) IsSupporting(elem task_element.Enum) bool {
	return [task_element.Len]bool{
		true, false, false,
	}[elem]
}

func (t *Task) ExecuteMain(args based.MainFuncArguments) error {

	t.reciever = make(chan string)
	t.close = make(chan error)

	callback := [task_element.Len]based.ImplCallback{nil}
	callback[task_element.StdIn] = based.MakeImplCallback(func(elem *entities.Element, buffer []byte) {

		name := elem.Fullname()
		format := t.Format[name]
		str := string(buffer)
		str = strings.ReplaceAll(str, "\n\r", "\n")
		str = strings.ReplaceAll(str, "\r\n", "\n")
		str = strings.ReplaceAll(str, "\r", "\n")
		str = strings.ReplaceAll(str, "\n", "\n"+format)

		fmt.Println(format + str + ForeReset + BackReset)

	}, based.DefaultOnFinal)

	go func() {

		args.ListenStart(callback)
		<-t.close
		args.Close(nil)

	}()

	return nil
}

func (t *Task) StopMain() {
	time.Sleep(50 * time.Millisecond)
	close(t.close)
}
