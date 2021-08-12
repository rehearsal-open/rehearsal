// entities/task.go
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
	"time"
)

// constant names task's kind, should be satisfied a entities.UserDefinedName regular expression
type TaskKind string

// constant names task's kind, should be satisfied a entities.UserDefinedName regular
type TaskDetailProperty int

type TaskDetail []interface{}

const (
	ExecuteCuiTask TaskKind = "executeCui"
	ConsoleTask    TaskKind = "console"
)

// defined configuration of rehearsal task, its lifespan
type Task struct {
	Name         DefinedName    // be satisfied a entities.DefinedName regular expression
	SyncInterval time.Duration  // i/o syncronization interval, default value is (*Default).Sync
	kind         TaskKind       // kind of task, select from entities.TaskKind enum
	SendTo       []TaskFullName // tasks' full name which i/o data send to, satisfied a entities.TaskFullName regular expression
	Details      TaskDetail     // detail configuration indivisually defined by task's kind
	BeginPhase   DefinedName    // name of phase which begins, must be equal to phase which appended it
	UntilPhase   DefinedName    // name of phase until which this task continue, default value is launching phase's name
}
