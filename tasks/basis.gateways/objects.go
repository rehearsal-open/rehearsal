// tasks/basis.gateways/objects.go
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

package basis

import (
	"errors"
	"io"

	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/element"
	"github.com/rehearsal-open/rehearsal/entities/run_state"
	"github.com/rehearsal-open/rehearsal/tasks/buffer"
)

type (
	// Defines main task's arguments
	MainFuncArguments interface {
		// get io.Writer using sender object
		Writer(element.TaskElement) (io.Writer, error)
		// call when main task's closing
		Close(err error)
		// call when main just after task's begining
		ListenStart() error
	}

	Task interface {
		MainState() run_state.RunningState
		ElementState(element element.TaskElement) run_state.RunningState
		BeginTask() error
		StopTask()
		WaitClosing()
		ReleaseResource()
		AppendReciever(sender element.TaskElement, reciever buffer.Reciever) error
		Reciever(element element.TaskElement) (buffer.Reciever, error)
	}

	RecieveCallback func(recieve buffer.Packet) error

	TaskImpl interface {
		IsSupporting(elem element.TaskElement) bool
		ExecuteMain(args MainFuncArguments) error
		StopMain()
		RecieverCallback(elem element.TaskElement) (RecieveCallback, error)
	}

	internalTask struct {
		impl       TaskImpl
		entity     *entities.Task
		numSupport int
		mainstate  run_state.RunningState
		state      [element.NumTaskElement]run_state.RunningState
		sender     [element.NumTaskElement]*buffer.Buffer
		reciever   [element.NumTaskElement]chan buffer.Packet
		closed     chan error
	}

	Element struct {
		*internalTask
		element.TaskElement
	}
)

var (
	ErrNotSupportingElement = errors.New("this task doesn't supports element.")
)
