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
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/frontend"
	"github.com/rehearsal-open/rehearsal/parser"
	"github.com/rehearsal-open/rehearsal/task/maker"
)

// Initialize rehearsal engine.
func (r *Rehearsal) Init(parser parser.Parser, envConfig parser.EnvConfig, maker *maker.Maker, frontend frontend.Frontend) error {

	// initialize instance
	(*r) = Rehearsal{
		frontend: frontend,
		lock:     &sync.Mutex{},
		entity:   &entities.Rehearsal{},
	}

	// parse entity
	if err := parser.Parse(envConfig, r.entity); err != nil {
		return errors.WithMessage(err, "rehearsal cannot parse")
	}

	// initialize task schedule array[1]
	r.tasks = []Task{}
	r.beginTasks = make([][]int, r.entity.NPhase)
	r.closeTasks = make([][]int, r.entity.NPhase)
	r.waitTasks = make([][]int, r.entity.NPhase)

	// Use to relation registering
	nameList := map[string]int{}

	// initialize task schedule array[2]
	for i := range r.beginTasks {
		r.beginTasks[i], r.closeTasks[i], r.waitTasks[i] = []int{}, []int{}, []int{}
	}

	// register entities.Task to engine.tasks
	if err := r.entity.ForeachTask(func(idx int, entity *entities.Task) error {

		// build task's executable instance
		if task, err := maker.MakeTask(entity); err != nil {
			return errors.WithMessage(err, "cannot make executable task instance")
		} else {

			// register task
			appended := len(r.tasks)
			nameList[task.Entity().Fullname()] = appended
			r.tasks = append(r.tasks, Task{
				Task:   task,
				entity: entity,
			})

			// append running schedules
			r.beginTasks[entity.LaunchAt] = append(r.beginTasks[entity.LaunchAt+1], appended)
			r.closeTasks[entity.CloseAt] = append(r.closeTasks[entity.CloseAt+1], appended)
			if entity.IsWait {
				r.waitTasks[entity.CloseAt] = append(r.waitTasks[entity.CloseAt+1], appended)
			}

			return nil
		}

	}); err != nil {
		return errors.WithMessage(err, "cannot register executable task instance")
	}

	// make frontend logger task
	if logger := r.frontend.LoggerTask(); logger != nil {
		appended := len(r.tasks)
		entity := logger.Entity()

		// set entity
		entity.Phasename, entity.Taskname = entities.SystemInitializePhase, "__frontend_logger"

		// register logger task
		nameList[entity.Fullname()] = appended
		r.tasks = append(r.tasks, Task{
			Task:   logger,
			entity: entity,
		})

		// register relation to logger task
		r.entity.ForeachTask(func(idx int, task *entities.Task) error {
			if task.Element[task_element.StdOut].WriteLog {
				task.AddRelation(task_element.StdOut, entity, task_element.StdIn)
			}
			if task.Element[task_element.StdErr].WriteLog {
				task.AddRelation(task_element.StdErr, entity, task_element.StdIn)
			}
			return nil
		})

		// append logger task's running schedule
		beginAt, _ := r.entity.Phase(entities.SystemInitializePhase)
		endAt, _ := r.entity.Phase(entities.SystemFinalizePhase)
		r.beginTasks[beginAt] = append(r.beginTasks[beginAt], appended)
		r.closeTasks[endAt] = append(r.closeTasks[endAt], appended)
	}

	// TODO: ADD SYSTEM TASK BEGINING AND CLOSING

	// append relation
	if err := r.entity.ForeachRelation(func(_ int, relation *entities.Relation) error {
		// get name
		senderName := relation.Sender.Fullname()
		recieverName := relation.Reciever.Fullname()

		// get task
		senderTask := r.tasks[nameList[senderName]].Task
		recieverTask := r.tasks[nameList[recieverName]].Task

		if reciever, err := recieverTask.Reciever(relation.ElementReciever); err != nil {
			return errors.WithMessage(err, "Cannot make relation from "+senderName+"'s "+relation.ElementSender.String()+" to "+recieverName+"'s "+relation.ElementReciever.String())
		} else if err := senderTask.AppendReciever(relation.ElementSender, reciever); err != nil {
			return errors.WithMessage(err, "Cannot make relation from "+senderName+"'s "+relation.ElementSender.String()+" to "+recieverName+"'s "+relation.ElementReciever.String())
		}
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
