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
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/task"
	"github.com/rehearsal-open/rehearsal/task/based"
	"github.com/rehearsal-open/rehearsal/task/buffer"
)

type (
	Task struct {
		based.Task
		Format   map[string]string
		reciever chan string
		close    chan error
		closed   chan error
	}
)

func MakeTask(entity *entities.Rehearsal) (task.Task, error) {

	maxNameLen := 0

	taskConf := entities.Task{}
	colorSet := [...]string{
		ForeRed, ForeGreen, ForeYellow, ForeBlue, ForeMagenta, ForeSyan,
	}

	task := Task{
		Format:   map[string]string{},
		reciever: make(chan string),
		close:    make(chan error),
		closed:   make(chan error),
	}

	task.Task = based.MakeBasis(&taskConf, &task)

	entity.Foreach(func(idx int, entity *entities.Task) error {
		if nameSize := len(entity.Fullname()); nameSize > maxNameLen {
			maxNameLen = nameSize
		}
		return nil
	})

	idx := 0

	entity.Foreach(func(_ int, entity *entities.Task) error {

		name := entity.Fullname()
		if entity.WriteLog {
			task.Format[name] = colorSet[idx%len(colorSet)] + BackReset + name + strings.Repeat(" ", maxNameLen-len(name)) + " : "
			idx++
		}

		return nil
	})

	return &task, nil
}

func (t *Task) IsSupporting(elem task_element.Enum) bool {
	return [task_element.Len]bool{
		true, false, false,
	}[elem]
}

func (t *Task) ExecuteMain(args based.MainFuncArguments) error {

	callback := [task_element.Len]based.Reciever{nil}
	callback[task_element.StdIn] = func(recieved *buffer.Packet) {

		sender, _ := recieved.Sender()
		defer recieved.Close()
		buffer := bytes.NewBuffer([]byte{})
		io.Copy(buffer, recieved)

		name := sender.Fullname()
		format := t.Format[name]
		str := buffer.String()
		str = strings.ReplaceAll(str, "\n\r", "\n")
		str = strings.ReplaceAll(str, "\r\n", "\n")
		str = strings.ReplaceAll(str, "\r", "\n")
		str = strings.ReplaceAll(str, "\n", "\n"+format)

		fmt.Println(format + str + ForeReset + BackReset)

	}

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
	// <-t.closed
}
