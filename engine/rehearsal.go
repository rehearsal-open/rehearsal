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
	"github.com/rehearsal-open/rehearsal/task"
	"github.com/rehearsal-open/rehearsal/task/maker"
)

// Initialize rehearsal engine.
func (r *Rehearsal) Init(parser parser.Parser, envConfig parser.EnvConfig, maker *maker.Maker, frontend frontend.Frontend) error {

	// initialize instance
	(*r) = Rehearsal{
		frontend: frontend,
		lock:     &sync.Mutex{},
		Entity:   &entities.Rehearsal{},
	}

	// parse entity
	if err := parser.Parse(envConfig, r.Entity); err != nil {
		return errors.WithMessage(err, "rehearsal cannot parse")
	}

	// initialize frontend
	frontend.Init(r.Entity)

	// initialize task schedule array[1]
	r.tasks = []Task{}
	r.beginTasks = make([][]int, r.Entity.NPhase)
	r.closeTasks = make([][]int, r.Entity.NPhase)
	r.waitTasks = make([][]int, r.Entity.NPhase)

	// Use to relation registering
	nameList := map[string]int{}

	// register task to rehearsal engine's instance
	registerTask := func(t task.Task) {

		appended := len(r.tasks)
		entity := t.Entity()

		r.tasks = append(r.tasks, Task{
			Task:     t,
			entity:   entity,
			frontend: frontend,
		})

		nameList[t.Entity().Fullname()] = appended

		r.beginTasks[entity.LaunchAt] = append(r.beginTasks[entity.LaunchAt], appended)
		r.closeTasks[entity.CloseAt] = append(r.closeTasks[entity.CloseAt], appended)
		if entity.IsWait {
			r.waitTasks[entity.CloseAt] = append(r.waitTasks[entity.CloseAt], appended)
		}

		// initialize element's parent(task configuration)
		for i, l := 0, task_element.Len; i < l; i++ {
			entity.Element[i].Parent = entity
			entity.Element[i].Kind = task_element.Enum(i)
		}

	}

	// initialize task schedule array[2]
	for i := range r.beginTasks {
		r.beginTasks[i], r.closeTasks[i], r.waitTasks[i] = []int{}, []int{}, []int{}
	}

	// register entities.Task to engine.tasks
	if err := r.Entity.ForeachTask(func(idx int, entity *entities.Task) error {

		// build task's executable instance
		if task, err := maker.MakeTask(entity); err != nil {
			return errors.WithMessage(err, "cannot make executable task instance")
		} else {
			// append running schedules
			registerTask(task)
			return nil
		}

	}); err != nil {
		return errors.WithMessage(err, "cannot register executable task instance")
	}

	// make frontend logger task
	if logger := r.frontend.LoggerTask(); logger != nil {

		entity := logger.Entity()

		// set entity
		entity.Phasename, entity.Taskname = entities.SystemInitializePhase, "__frontend_logger"
		entity.LaunchAt, _ = r.Entity.Phase(entities.SystemInitializePhase)
		entity.CloseAt, _ = r.Entity.Phase(entities.SystemFinalizePhase)

		// initialize element's parent(task configuration)
		for i, l := 0, task_element.Len; i < l; i++ {
			entity.Element[i].Parent = entity
			entity.Element[i].Kind = task_element.Enum(i)
		}

		// register relation to logger task
		r.Entity.ForeachTask(func(idx int, task *entities.Task) error {
			if task.Element[task_element.StdOut].WriteLog {
				task.AddRelation(task_element.StdOut, entity, task_element.StdIn)
			}
			if task.Element[task_element.StdErr].WriteLog {
				task.AddRelation(task_element.StdErr, entity, task_element.StdIn)
			}
			return nil
		})

		registerTask(logger)
	}

	// TODO: ADD SYSTEM TASK BEGINING AND CLOSING

	// append relation
	if err := r.Entity.ForeachRelation(func(_ int, relation *entities.Relation) error {
		// get name
		senderName := relation.Sender.Fullname()
		recieverName := relation.Reciever.Fullname()

		// get task
		senderTask := r.tasks[nameList[senderName]].Task
		recieverTask := r.tasks[nameList[recieverName]].Task

		if err := senderTask.Connect(relation.ElementSender, relation.ElementReciever, recieverTask); err != nil {
			return errors.WithMessage(err, "Cannot make relation from "+senderName+"'s "+relation.ElementSender.String()+" to "+recieverName+"'s "+relation.ElementReciever.String())
		}

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
