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

import (
	"regexp"

	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/streamwest-1629/textfilter"
)

type (
	// Defines configuration of rehearsal excuting and each task  default configuration's value.
	Rehearsal struct {
		Version       float64        `map-to:"version!"` // Supported 1.202109 only.
		DefaultDir    string         `map-to:"dir"`      // Using cui task as default execute directory.
		tasks         []*Task        // Tasks.
		tasknameList  map[string]int // Pair each task's index and its name in tasks.
		phasenameList map[string]int // Pair each phase's index and its name in tasks.
		NPhase        int            // The number of phase without system initialize/finalize phase.
	}

	// Defines configuration of rehearsal task, its lifespan.
	Task struct {
		Phasename string // Name of task's began phase.
		Taskname  string `map-to:"name!"` // Name of task.
		// Kind name of task. Supported name is defined by /task/maker's structure.
		Kind string `map-to:"kind!"`
		// Task's launching phase number. Defined by /engine's structure.
		LaunchAt int
		// Task's closing phase number. Defined by /engine's structure.
		CloseAt int
		// Whether task's natural closing.
		// When it is true, engine closes task after task's natural closing.
		// If not, engine closes without waiting for task's natural closing.
		//
		// Default is true.
		IsWait bool
		// Task's detail interface.
		// Instance's type is differented by task's kind.
		Detail TaskDetail
		// Configurations task's element.
		Element [task_element.Len]Element
	}

	// Defines configuration of rehearsal task's element.
	Element struct {
		// Whether send task's
		WriteLog bool `map-to:"write-log"`
		Sendto   []Relation
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
	Relation struct {
		Sender          *Task
		Reciever        *Task
		ElementSender   task_element.Enum
		ElementReciever task_element.Enum
	}
)

var (
	UserShortNameRegexp  = textfilter.RegexpMatches(`^[a-zA-Z][a-zA-Z0-9_]*$`)
	UserFullNameRegexp   = textfilter.RegexpMatches(`^[a-zA-Z][a-zA-Z0-9_]*::[a-zA-Z][a-zA-Z0-9_]*$`)
	fullNameParserRegexp = [...]*regexp.Regexp{
		regexp.MustCompile(`^(?P<phase>((__)?[a-zA-Z][a-zA-Z0-9_]*))::(?P<task>((__)?[a-zA-Z][a-zA-Z0-9_]*))#(?P<element>(stdin|stdout|stderr))$`),
		regexp.MustCompile(`^(?P<phase>((__)?[a-zA-Z][a-zA-Z0-9_]*))::(?P<task>((__)?[a-zA-Z][a-zA-Z0-9_]*))$`),
		regexp.MustCompile(`^(?P<task>((__)?[a-zA-Z][a-zA-Z0-9_]*))#(?P<element>(stdin|stdout|stderr))$`),
		regexp.MustCompile(`^(?P<task>((__)?[a-zA-Z][a-zA-Z0-9_]*))$`),
	}
)

const (
	SystemInitializePhase = "__init"
	SystemFinalizePhase   = "__finl"
	UserFinalizePhase     = "__endless"
)
