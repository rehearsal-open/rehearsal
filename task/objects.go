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
	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_state"
	"github.com/rehearsal-open/rehearsal/task/buffer"
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
		// Gets selected element's running state.
		// When try to get unsupported elements, return value is undefined.
		//
		// See each element's details also: github.com/rehearsal-open/rehearsal/entities/enum/task_state
		ElementState(element task_element.Enum) task_state.Enum
		// Begin task. Internal functions, execute goroutine main task and some supporting elements.
		BeginTask() error // begin main task
		// Force to stop task.
		StopTask() // stop reciever and main task
		// Wait for natural stopping task or force stopping task.
		WaitClosing()
		// Release memory and any handler.
		ReleaseResource() // delete buffer and so on
		// Append reciever to selected sender element.
		AppendReciever(sender task_element.Enum, reciever buffer.Reciever) error
		// Get reciever selected by task element.
		Reciever(element task_element.Enum) (buffer.Reciever, error)
	}

	// Definates functions to make task's detail instance.
	DetailCreator interface {
		// Assign task's detail with configuration.
		AssignTaskDetail(kind string, entity *entities.Task, mapped map[string]interface{}) error
	}

	// Definates functions to make task's instance.
	TaskCreator interface {
		// Create new task instance.
		MakeTask(kind string, entity *entities.Task) (Task, error)
	}
)
