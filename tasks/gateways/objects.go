// tasks/gateways/objects.go
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

package gateways

import (
	"github.com/rehearsal-open/rehearsal/entities/element"
	"github.com/rehearsal-open/rehearsal/entities/run_state"
	"github.com/rehearsal-open/rehearsal/tasks/buffer"
)

type (
	Task interface {
		IsSupporting(element element.TaskElement) bool
		MainState() run_state.RunningState
		ElementState(element element.TaskElement) run_state.RunningState
		BeginTask() error // begin main task
		StopTask()        // stop reciever and main task
		ReleaseResource() // delete buffer and so on
		AppendReciever(sender element.TaskElement, reciever buffer.Reciever) error
		Reciever(element element.TaskElement) (buffer.Reciever, error)
	}
)
