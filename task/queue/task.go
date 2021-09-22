// task/queue/task.go
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

package queue

import (
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/task"
)

func Connect(fromOut task.Task, fromElem task_element.Enum, toIn task.Task, toElem task_element.Enum) error {

	// check fromOut and toIn argument are valid supported value or not.
	if from, ok := fromOut.(Task); !ok {
		panic("fromOut argument is invalid value type(un-supported queue)")
	} else if to, ok := toIn.(Task); !ok {
		panic("toIn argument is invalid value type(un-supported queue)")
	} else {

		// get element
		if out := from.GetOutput(fromElem); out == nil {
			return task.ErrNotSupportingElement
		} else if in := to.GetInput(toElem); in == nil {
			return task.ErrNotSupportingElement
		} else {

			// connect
			out.AppendInput(MakeWriter(in))
			return nil

		}
	}
}
