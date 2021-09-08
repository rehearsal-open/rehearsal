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
	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/task"
	"github.com/rehearsal-open/rehearsal/task/based"
	"github.com/rehearsal-open/rehearsal/task/queue"
	"github.com/rehearsal-open/rehearsal/task/wrapper"
	"github.com/rehearsal-open/rehearsal/task/wrapper/listen"
)

func (parallel *ElemParallel) AppendElem(fromElem entities.Element, insert wrapper.Filter, finallyElem task_element.Enum) {
	name := fromElem.Fullname()
	if _, exist := parallel.parallelWriter[name]; exist {
		panic("already registered element")
	} else {
		if sendto := parallel.GetInput(finallyElem); sendto == nil {
			panic(task.ErrNotSupportingElement.Error())
		} else {
			insert.Output().AppendWriteTo(sendto)
			parallel.parallelWriter[name] = insert
		}

	}
}

func (parallel *ElemParallel) IsSupporting(elem task_element.Enum) bool {
	return [task_element.Len]bool{
		true, parallel.finallyTask.IsSupporting(1), parallel.finallyTask.IsSupporting(2),
	}[elem]
}

func (parallel *ElemParallel) ExecuteMain(args based.MainFuncArguments) error {

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

func (parallel *ElemParallel) StopMain() {
	close(parallel.close)
}

func (parallel *ElemParallel) GetOutput(elem task_element.Enum) *queue.Senders {
	return parallel.finallyTask.GetOutput(elem)
}
