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

import "strconv"

func (r *Rehearsal) Execute() error {

	defer r.releaseResources()

	for i := range r.beginTasks {
		beginTasks, waitTasks, closeTasks := r.beginTasks[i], r.waitTasks[i], r.closeTasks[i]

		r.frontend.Log(0, "start phase ("+strconv.Itoa(i)+" / "+strconv.Itoa(len(r.beginTasks))+")")

		for _, val := range beginTasks {

			if err := r.tasks[val].BeginTask(); err != nil {
				return err
			}

			r.frontend.Log(0, "running start: "+r.tasks[val].entity.Taskname)
		}

		for _, val := range waitTasks {
			r.tasks[val].WaitClosing()
		}

		for _, val := range closeTasks {
			r.tasks[val].StopTask()
			r.frontend.Log(0, "task closed ("+r.tasks[val].entity.Fullname()+")")
		}
	}
	return nil
}

func (r *Rehearsal) releaseResources() {
	for i := range r.tasks {
		r.tasks[i].ReleaseResource()
	}
}