// task/objects.go
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

package task

import (
	"errors"

	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_state"
)

type (
	// Definates functions to be task.
	//
	// Tasks should be embedding basis.Task interface.
	// Almosts of these function included by it: github.com/rehearsal-open/rehearsal/tasks/based
	Task interface {
		// Gets task's configuration in entities.
		Entity() *entities.Task
		// Gets whether element is supported by task.
		IsSupporting(element task_element.Enum) bool
		// Gets selected element's running state.
		//
		// See each element's details also: github.com/rehearsal-open/rehearsal/entities/enum/task_state
		MainState() task_state.Enum
		// Begin task. Internal functions, execute goroutine main task and some supporting elements.
		BeginTask() error // begin main task
		// Force to stop task.
		StopTask() // stop reciever and main task
		// Wait for natural stopping task or force stopping task.
		WaitClosing()
		// Release memory and any handler.
		ReleaseResource() // delete buffer and so on
		// Append reciever to selected sender element.
		Connect(senderElem task_element.Enum, recieverElem task_element.Enum, reciever Task) error
	}
)

var (
	ErrNotSupportingElement = errors.New("this task doesn't supports element")
	ErrAlreadyRun           = errors.New("task has already run")
	ErrAlreadyClosed        = errors.New("task has already closed")
)
