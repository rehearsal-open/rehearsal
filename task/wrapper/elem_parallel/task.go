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
	"github.com/rehearsal-open/rehearsal/task/based"
)

func (parallel *__task) IsSupporting(elem task_element.Enum) bool {
	return [task_element.Len]bool{
		true, parallel.finallyTask.IsSupporting(1), parallel.finallyTask.IsSupporting(2),
	}[elem]
}

func (parallel *__task) ExecuteMain(args based.MainFuncArguments) error {

	callback := [task_element.Len]based.ImplCallback{nil}
	callback[task_element.StdIn] = based.MakeImplCallback(func(elem *entities.Element, b []byte) {
		name := elem.Fullname()
		if writer, exist := parallel.parallelWriter[name]; exist {
			writer.Write(elem, b)
		}
	}, func() {
		// TODO: WRITER.CLOSE()
	})
}
