// entities/rehearsal.go
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

import "github.com/pkg/errors"

func (r *Rehearsal) AddTask(task *Task) {
	at, name := len(r.tasks), task.Fullname()
	r.tasks = append(r.tasks, task)
	if r.nameList == nil {
		r.nameList = make(map[string]int)
	}
	r.nameList[name] = at
}

// For-each loop appended task's entity
func (r *Rehearsal) Foreach(action func(idx int, task *Task) error) error {

	for i, task := range r.tasks {
		if err := action(i, task); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// Return number of including task.
func (r *Rehearsal) LenTask() int { return len(r.tasks) }

// Return task instance selected by identifier name.
func (r *Rehearsal) Task(fullname string) (*Task, error) {
	if idx, ok := r.nameList[fullname]; !ok {
		return nil, ErrCannotFoundProperty("task", fullname)
	} else {
		return r.tasks[idx], nil
	}
}
