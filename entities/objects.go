// entities/objects.go
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

package entities

import "github.com/rehearsal-open/rehearsal/entities/enum/task_element"

type (
	// Defines configuration of rehearsal excuting and each task  default configuration's value.
	Rehearsal struct {
		Version  float64 `map-to:"version!"`
		tasks    []*Task
		nameList map[string]int
	}

	// Defines configuration of rehearsal task, its lifespan.
	Task struct {
		Phasename string
		Taskname  string `map-to:"name!"`
		Kind      string `map-to:"kind!"`
		LaunchAt  int
		Detail    TaskDetail
		sendto    []Reciever
	}

	// Defines functions whose task's detail structure must be statisfied as task's detail structure.
	TaskDetail interface {
		// Validate member value.@
		// If it is able to fix them, should do that.
		CheckFormat() error
		// Convert from TaskDetail to string.
		String() string
	}

	// Defines relation bitween task and task.
	Reciever struct {
		Reciever *Task
		Element  task_element.Enum
	}
)
