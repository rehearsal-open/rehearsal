// task/wrapper/elem_parallel/make.go
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
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/task/based"
	"github.com/rehearsal-open/rehearsal/task/queue"
	"github.com/rehearsal-open/rehearsal/task/wrapper"
)

func Make(finally queue.Task) *ElemParallel {

	result := ElemParallel{
		isSupport:      [3]bool{true, false, false},
		parallelWriter: make(map[string]wrapper.Filter),
		finallyTask:    finally,
		close:          make(chan error),
	}

	result.Task = based.MakeBasis(finally.Entity(), &result)
	result.isSupport = [...]bool{true, finally.IsSupporting(task_element.Enum(1)), finally.IsSupporting(task_element.Enum(2))}

	return &result
}
