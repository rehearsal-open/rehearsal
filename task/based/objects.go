// task/based/objects.go
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

package based

import (
	"errors"
	"io"
	"sync"

	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_state"
	"github.com/rehearsal-open/rehearsal/task/buffer"
)

type (
	implCallback struct {
		reciever func(recieved *buffer.Packet)
		onfinal  func()
	}

	// Callback functions called when task recieves packet.
	ImplCallback interface {
		Recieve(recieved *buffer.Packet)
		OnFinal()
	}

	// Defines main task's argumetns to call basis functions.
	// These functions must to be called.
	MainFuncArguments interface {
		// Gets io.Writer using sender object.
		Writer(task_element.Enum) (io.Writer, error)
		// Call basis that main task is closed. Main task must call when main task's closing.
		Close(err error)
		// Listen reciever element and begin to manage packet.
		// Main task must call when main just after task's begining.
		ListenStart(callback [task_element.Len]ImplCallback)
	}

	// Defines functions basis task including.
	// Implemented tasks should use this interface as embedded interface.
	//
	// These functions are partially satisfied gateways.Task interface.
	// See functions' detail: github.com/rehearsal-open/rehearsal/enum/task_element
	Task interface {
		// Gets task's configuration in entities.
		Entity() *entities.Task
		// Gets main task's running state.
		MainState() task_state.Enum
		// Gets selected element's running state.
		ElementState(element task_element.Enum) task_state.Enum
		// Begin main task.
		BeginTask() error
		// Stop main task.
		StopTask()
		// Wait for closing main task.
		WaitClosing()
		// Release memory and any handler.
		ReleaseResource()
		// Append reciever to selected sender element.
		AppendReciever(sender task_element.Enum, reciever buffer.SendToBased) error
		// Get reciever selected by task element.
		Reciever(element task_element.Enum) (buffer.SendToBased, error)
	}

	// Defines functions which implemented tasks are satisfied.
	TaskImpl interface {
		// Gets whether element is supported by task.
		IsSupporting(task_element.Enum) bool
		// Run main task.
		// MainFuncArguments interface must be used to call basis function internal main task's function.
		// Main task's function must call its member methods.
		// See MainFuncArguments' details.
		ExecuteMain(args MainFuncArguments) error
		// Stop main task
		StopMain()
	}

	internalTask struct {
		impl      TaskImpl
		entity    *entities.Task
		mainstate task_state.Enum
		elements  [task_element.Len]taskElement
		lock      sync.Mutex
		closed    chan error
	}

	taskElement struct {
		*internalTask
		lock       sync.Mutex
		element    task_element.Enum
		state      task_state.Enum
		sender     *buffer.Buffer
		reciever   chan buffer.Packet
		packets    []buffer.Packet
		packetPos  int
		packetLock *sync.Mutex
	}
)

var (
	ErrNotSupportingElement = errors.New("this task doesn't supports element")
	ErrAlreadyRun           = errors.New("task has already run")
	ErrAlreadyClosed        = errors.New("task has already closed")
)
