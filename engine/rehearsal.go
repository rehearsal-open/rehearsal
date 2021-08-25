// engine/rehearsal.go
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
	"sync"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/frontend"
	"github.com/rehearsal-open/rehearsal/parser"
	"github.com/rehearsal-open/rehearsal/task/maker"
)

func (r *Rehearsal) Reset(parser parser.Parser, maker maker.Maker, frontend frontend.Frontend) error {

	(*r) = Rehearsal{
		frontend: frontend,
		lock:     &sync.Mutex{},
	}

	var nEntity int
	nameList := map[string]int{}

	// parser entity
	if entity, err := parser.Parse(); err != nil {
		return errors.WithStack(err)
	} else {
		r.entity = entity
		nEntity = entity.NPhase + 2
	}

	r.tasks = make([]Task, 0, r.entity.LenTask()*2)
	r.beginTasks = make([][]int, nEntity)
	r.closeTasks = make([][]int, nEntity)
	r.waitTasks = make([][]int, nEntity)

	for i := range r.beginTasks {
		r.beginTasks[i], r.closeTasks[i], r.waitTasks[i] = []int{}, []int{}, []int{}
	}

	// register task
	if err := r.entity.Foreach(func(idx int, entity *entities.Task) error {

		// build task
		if task, err := maker.MakeTask(entity); err != nil {
			return errors.WithStack(err)
		} else {

			appended := len(r.tasks)
			nameList[entity.Fullname()] = appended
			r.tasks = append(r.tasks, Task{
				Task:   task,
				entity: entity,
			})

			// append running schedules
			r.beginTasks[entity.LaunchAt+1] = append(r.beginTasks[entity.LaunchAt+1], appended)
			r.closeTasks[entity.CloseAt+1] = append(r.closeTasks[entity.CloseAt+1], appended)
			if entity.IsWait {
				r.waitTasks[entity.CloseAt+1] = append(r.waitTasks[entity.CloseAt+1], appended)
			}

			return nil
		}

	}); err != nil {
		return errors.WithStack(err)
	}

	// TODO: ADD SYSTEM TASK BEGINING AND CLOSING

	// append relation
	for i := range r.tasks {
		task := r.tasks[i]

		// relation foreach-loop
		if err := task.Entity().RelationForeach(func(idx int, relation *entities.Reciever) error {

			// TODO: SYSTEM RELATION BITWEEN TASK AND TASK
			recieverTask := r.tasks[nameList[relation.Reciever.Fullname()]]
			if reciever, err := recieverTask.Reciever(relation.ElementReciever); err != nil {
				return errors.WithStack(err)
			} else if err := task.AppendReciever(relation.ElementSender, reciever); err != nil {
				return errors.WithStack(err)
			}
			return nil
		}); err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}
