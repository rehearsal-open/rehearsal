// entities/task.go
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

// Show task's fullname, including task's phasename and taskname.
func (t *Task) Fullname() string {
	return t.Phasename + "::" + t.Taskname
}

// Show task's fullname and element.
func (t *Task) FullnameWithElem(elem task_element.Enum) string {
	return t.Fullname() + "(" + elem.String() + ")"
}

func (t *Task) AddRelation(output task_element.Enum, recieverTask *Task, recieverElem task_element.Enum) {
	t.Element[output].Sendto = append(t.Element[output].Sendto, Relation{
		Sender:          t,
		Reciever:        recieverTask,
		ElementSender:   output,
		ElementReciever: recieverElem,
	})
}

// Return number of including relation.
func (t *Task) LenRelation() int {
	ans := 0
	for i := range t.Element {
		ans += len(t.Element[i].Sendto)
	}
	return ans
}

// Foreach-loop in relation which send from this task.
func (t *Task) RelationForeach(action func(idx int, output task_element.Enum, relation *Relation) error) error {

	idx := 0

	for i := range t.Element {
		for j := range t.Element[i].Sendto {
			if err := action(idx, task_element.Enum(i), &t.Element[i].Sendto[j]); err != nil {
				return errors.WithStack(err)
			} else {
				idx++
			}
		}
	}
	return nil
}
