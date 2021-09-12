// engine/task.go
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

package engine

import (
	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_state"
)

// Begin task and prints log.
func (task *Task) BeginTask() error {
	if err := task.Task.BeginTask(); err != nil {
		return errors.WithStack(err)
	}

	task.frontend.Log(0, task.entity.Fullname()+" was started.")
	return nil
}

// Stop task and prints log.
func (task *Task) StopTask() {
	switch task.Task.MainState() {
	case task_state.Closed, task_state.Finalized:
		task.frontend.Log(0, task.entity.Fullname()+" has been already stopped.")
	case task_state.Running:
		task.Task.StopTask()
		task.frontend.Log(0, task.entity.Fullname()+" was stopped.")
	case task_state.Waiting:
		task.Task.StopTask()
	}
}

// Release resources and prints log.
func (task *Task) ReleaseResource() {
	switch task.Task.MainState() {
	case task_state.Closed:
		task.Task.ReleaseResource()
	case task_state.Running:
		task.Task.ReleaseResource()
		task.frontend.Log(0, task.entity.Fullname()+" were forced to be stopped.")
	case task_state.Waiting:
		task.Task.ReleaseResource()
	}
}
