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

import (
	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
)

// Add task.
// If duplicated name, it occers error.
func (r *Rehearsal) AddTask(task *Task) error {
	at, name := len(r.tasks), task.Fullname()
	if r.nameList == nil {
		r.nameList = make(map[string]int)
	} else if _, exist := r.nameList[name]; exist {
		return errors.New("duplicated task's name in same phase: " + name)
	}
	r.tasks = append(r.tasks, task)
	r.nameList[name] = at
	return nil
}

// For-each loop appended task's entity.
func (r *Rehearsal) ForeachTask(action func(idx int, task *Task) error) error {

	for i, task := range r.tasks {
		if err := action(i, task); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// For-each loop appended all task's relation.
func (r *Rehearsal) ForeachRelation(action func(idx int, reciever *Relation) error) error {

	idx := 0

	for i := range r.tasks { // foreach tasks in rehearsal config
		for j := range r.tasks[i].Element { // foreach element in task
			for k := range r.tasks[i].Element[j].Sendto { // for each relation in task's element
				relation := r.tasks[i].Element[j].Sendto[k]
				if err := action(idx, &relation); err != nil {
					return errors.WithStack(err)
				}
				idx++
			}
		}
	}
	return nil
}

// Parse to fullname to phasename, taskname, and element.
// When successfully, err is nil, if not, results without err is undefined.
func ParseTaskElem(fullname string, defaultPhase string, defaultElem task_element.Enum) (phaseName string, taskName string, element task_element.Enum, err error) {

	const (
		PhaseTemp = `$phase`
		TaskTemp  = `$task`
		ElemTemp  = `$element`
	)

	for i := range fullNameParserRegexp {
		if match := fullNameParserRegexp[i].FindStringSubmatchIndex(fullname); match != nil {
			reg := fullNameParserRegexp[i]
			switch i {
			case 0: // fullname including element
				// set parse result
				phaseName = string(reg.ExpandString([]byte{}, PhaseTemp, fullname, match))
				taskName = string(reg.ExpandString([]byte{}, TaskTemp, fullname, match))
				elemName := string(reg.ExpandString([]byte{}, ElemTemp, fullname, match))
				element = task_element.Parse(elemName)

				// check varidated element's name successfully
				if element == task_element.Unknown {
					err = errors.WithMessage(task_element.ErrUnknownElement, "cannot parse expression because of its element: "+elemName)
				} else {
					err = nil
				}
				return

			case 1: // fullname
				// set parse result
				phaseName = string(reg.ExpandString([]byte{}, PhaseTemp, fullname, match))
				taskName = string(reg.ExpandString([]byte{}, TaskTemp, fullname, match))
				element, err = defaultElem, nil
				return

			case 2: // shortname and element
				// set parse result
				taskName = string(reg.ExpandString([]byte{}, TaskTemp, fullname, match))
				elemName := string(reg.ExpandString([]byte{}, ElemTemp, fullname, match))
				phaseName = defaultPhase
				element = task_element.Parse(elemName)

				// check varidated element's name successfully
				if element == task_element.Unknown {
					err = errors.WithMessage(task_element.ErrUnknownElement, "cannot parse expression because of its element: "+elemName)
				}
				return

			case 3: // shortname
				taskName = string(reg.ExpandString([]byte{}, TaskTemp, fullname, match))
				phaseName, element, err = defaultPhase, defaultElem, nil
				return
			}
		}
	}

	return "", "", task_element.Unknown, errors.New("cannot parse expression because it is invalid format: " + fullname)
}
