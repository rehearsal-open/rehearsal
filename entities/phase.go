// entities/phase.go
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

import "errors"

func (p *Phase) Name() string { return p.name }

// Make new Task instance and append it to this instance.
func (p *Phase) appendTask(name string, kind string, detail TaskDetail) (task *Task, err error) {

	successful := func(_ string) (done bool, err error) {

		// when filtered successfully

		// make task
		task = &Task{
			name:         name,
			Kind:         kind,
			SyncInterval: p.SyncInterval,
			detail:       detail,
			Encoding:     "utf8",
		}

		// append task
		idx := len(p.tasks)
		p.tasks = append(p.tasks, task)
		p.nameList[name] = idx

		return true, nil
	}

	// validate task's name
	if _, err := p.taskFilter.Add(name,
		func(_ string) (done bool, err error) {
			return taskkindFilter.Add(kind, successful)
		}); err != nil {

		// invalid name
		return nil, err
	}

	return task, nil

}

// Return task object selected by order.
func (p *Phase) TaskAt(idx int) (*Task, error) {
	if idx > -1 && idx < len(p.tasks) {
		return p.tasks[idx], nil
	} else {
		return nil, errors.New("task's order is overed")
	}
}

// Return number of including tasks.
func (p *Phase) TaskLen() int {
	return len(p.tasks)
}

// Return task object selected by name.
func (p *Phase) Task(name string) (*Task, error) {
	if idx, ok := p.nameList[name]; !ok {
		return nil, errors.New("cannot found phase's name")
	} else {
		return p.tasks[idx], nil
	}
}
