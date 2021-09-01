// engine/v0/execute.go
// Copyright (C) 2021  Kasai Koji

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

package v0

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/task"
)

func (e *RehearsalEngine) Run() error {

	type exitTask struct {
		err  error
		task task.Task
	}

	exit := make(chan exitTask)
	nTask := len(e.tasks)
	iTask := 0

	for name, t := range e.tasks {
		e.logger.SystemPrint(fmt.Sprint("initialize task (", 1+iTask, "/", nTask, " : ", name, ")..."))
		if err := t.RunInit(); err != nil {
			return errors.WithStack(err)
		}
		iTask++
	}

	iTask = 0

	for _, t := range e.tasks {

		e.logger.SystemPrint(fmt.Sprint("running start(", iTask+1, "/", nTask, " : ", t.Config().Name, ")..."))
		go func(t task.Task, exit chan exitTask) {
			exit <- exitTask{
				err:  t.RunWait(),
				task: t,
			}
		}(t, exit)
		iTask++
	}

	defer func() {
		for _, t := range e.tasks {
			e.logger.SystemPrint(fmt.Sprint("finalize: " + t.Config().Name))
			t.Finalize()
			e.logger.SystemPrint("finished")
		}
	}()

	iTask = 0

	for iTask < nTask {

		select {
		case exited := <-exit:
			if exited.err != nil {
				e.logger.SystemPrint(fmt.Sprint("error occered at ", exited.task.Config().Name, ": ", exited.err.Error()))
			}
			e.logger.SystemPrint(fmt.Sprint("running end(", iTask+1, "/", nTask, " : ", exited.task.Config().Name, ")"))
			iTask++
		default:
			time.Sleep(time.Duration(e.config.SyncMs))
		}
	}

	return nil
}

func (e *RehearsalEngine) Kill() {
	// todo: write definition
}
