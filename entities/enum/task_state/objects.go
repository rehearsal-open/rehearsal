// entities/enum/task_state/objects.go
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

package task_state

// Definates task's running state.
type Enum int

const (
	// Task is waiting, in other words, this state is previous running.
	// Although in this state, task can recieves packets.
	// But should not get it in reciever task element.
	//
	// If task is called function to close task, this task don't running.
	Waiting Enum = iota
	// Task is running.
	// Reciever task elements and sender task elements are running too.
	//
	// If task is called function duplicated begining to run task,
	// Task is once running when only first calling.
	Running
	// Task is closed, but task don't release memory or other handler yet.
	//
	// If task is called function to run task, returns error showing already closed.
	// Then if task is called function to duplicated close task,
	// function should do nothing.
	Closed
	// Task is finalized, it means task release memory or other handler.
	//
	// If task is called function to run task, returns error showing already closed.
	// Then if task is called function to close or duplicated to finalize task,
	// function should do nothing.
	Finalized
)
