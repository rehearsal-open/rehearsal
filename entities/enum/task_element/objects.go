// entities/enum/task_element/objects.go
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

package task_element

import "errors"

// Definates kind of task's elements.
// An element works for managing task's monitoring.
// TaskElement is without main task.
type Enum int

const (
	// Unknown value.
	Unknown Enum = -1
	// Standard input element. This is reciever task's element.
	StdIn Enum = iota
	// Standard output element. This is sender task's element.
	StdOut
	// Standard Error output element. This is sender task's element.
	StdErr
	numTaskElement
	// The number what kinds of task elements.
	Len = int(numTaskElement)
)

// Expresses from task element's enum to string.
func (t Enum) String() string {
	return [Len]string{
		"standard input",
		"standard output",
		"standard error output",
	}[t]
}

func Parse(str string) Enum {
	switch str {
	case "stdin":
		return StdIn
	case "stdout":
		return StdOut
	case "stderr":
		return StdErr
	default:
		return Unknown
	}
}

var (
	ErrUnknownElement = errors.New("this is unknown element")
)
