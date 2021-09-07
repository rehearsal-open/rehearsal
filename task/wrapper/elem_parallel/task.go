// task/wrapper/elem_parallel/task.go
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

package elem_parallel

import (
	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/task"
	"github.com/rehearsal-open/rehearsal/task/based"
	"github.com/rehearsal-open/rehearsal/task/queue"
	"github.com/rehearsal-open/rehearsal/task/wrapper"
	"github.com/rehearsal-open/rehearsal/task/wrapper/listen"
)

func (parallel *__task) AppendElem(fromElem entities.Element, insert task.Task, inElem task_element.Enum, outElem task_element.Enum) error {
	name := fromElem.Fullname()
	if _, exist := parallel.parallelWriter[name]; exist {
		panic("already registered element")
	} else {
		if insertIn := wrapper.GetQueueAccess(insert).GetInput(inElem); insertIn == nil {
			return errors.WithStack(task.ErrNotSupportingElement)
		} else {
			parallel.parallelWriter[name] = queue.MakeWriter(insertIn)
		}
		return errors.WithStack(queue.Connect(insert, outElem, parallel.finallyTask, task_element.StdIn /* TODO?: SHOULD BE TO MEMBER */))
	}
}

func (parallel *__task) IsSupporting(elem task_element.Enum) bool {
	return [task_element.Len]bool{
		true, parallel.finallyTask.IsSupporting(1), parallel.finallyTask.IsSupporting(2),
	}[elem]
}

func (parallel *__task) ExecuteMain(args based.MainFuncArguments) error {

	listen.ListenElemBytes(parallel, task_element.StdIn, func(elem *entities.Element, b []byte) {
		name := elem.Fullname()
		if writer, exist := parallel.parallelWriter[name]; exist {
			writer.Write(elem, b)
		}
	}, func() {
		for key := range parallel.parallelWriter {
			parallel.parallelWriter[key].Close()
		}
	})

	go func() {
		<-parallel.close
		args.Close(nil)
	}()

	return nil
}

func (parallel *__task) StopMain() {
	close(parallel.close)
}
