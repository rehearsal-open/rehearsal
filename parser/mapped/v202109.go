// parser/mapped/v202109.go
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

package mapped

import (
	"math"
	"sort"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/streamwest-1629/textfilter"
)

// Parse object with rehearsal
func (p *Parser) Parse202109(r *Rehearsal) error {

	// insert initialize and finalize phases to counter.
	phases := map[string]*Phase{}
	initPhases := []*Phase{
		{
			Name:  entities.SystemInitializePhase,
			Tasks: []Task{},
		},
	}
	finalPhases := []*Phase{
		{
			Name: entities.UserFinalizePhase,
			Tasks: []Task{},
		}
		{
			Name:  entities.SystemFinalizePhase,
			Tasks: []Task{},
		},
	}
	r.Phases = append(initPhases, r.Phases...)
	r.Phases = append(r.Phases, finalPhases...)
	r.Rehearsal.NPhase = len(r.Phases)
	nInitPhase, nFinalPhase := len(initPhases), len(finalPhases)

	// use to check phase's name is valid or not
	phaseFilter := textfilter.Multiple([]textfilter.Filter{
		textfilter.Identifier(),
		entities.UserShortNameRegexp,
	})

	// use to check task'name is valid or not
	taskFilter := textfilter.Multiple{
		textfilter.Identifier(),
		entities.UserFullNameRegexp,
	}

	// check phase name and set index
	for iPhase := range r.Phases {
		if iPhase >= nInitPhase && iPhase < r.Rehearsal.NPhase-nFinalPhase {
			phase := r.Phases[iPhase] // caching

			// varidate
			if err := textfilter.RegisterFiltering(phaseFilter, phase.Name, func() error {
				if phase.Index < 1 { // order check
					phase.Index = math.MaxInt32 // default value
				}
				return nil
			}); err != nil {
				return errors.WithMessage(err, errMsgBase+"Phase[i] property: ("+phase.Name+")")
			}

		}
	}

	// sort phase
	sort.Sort(phaseByIndex(r.Phases[nInitPhase : r.NPhase-nFinalPhase]))

	// assign sorted result to phase
	for iPhase := range r.Phases {
		phase := r.Phases[iPhase]

		phase.Index = iPhase

		// register phase by name
		phases[phase.Name] = phase

		r.SetPhase(phase.Name, iPhase)
	}

	// assign task's details
	for iPhase := range r.Phases {
		phase := r.Phases[iPhase]

		for iTask := range phase.Tasks {
			task := &phase.Tasks[iTask]
			entity := task.Task

			task.Phasename = phase.Name

			// check task's name is valid or not
			if err := textfilter.RegisterFiltering(taskFilter, entity.Fullname(), func() error {

				// set task's launch, wait, close default value
				entity.LaunchAt, entity.CloseAt = iPhase, iPhase

				// set task detail data
				if err := p.DetailMaker.MakeDetail(r.Rehearsal, task.Clone, task.Task); err != nil {
					return errors.WithMessage(err, errMsgBase+task.Fullname()+"'s details")
				}

				r.Rehearsal.AddTask(task.Task)

				// validate until's property(phase name)
				if task.UntilPhase != nil {
					if until, err := r.Rehearsal.Phase(*task.UntilPhase); err != nil {
						return errors.WithMessage(err, errMsgBase+task.Fullname()+"'s until property's name(it is phase name)")
					} else if entity.LaunchAt > until {
						return errors.WithMessage(err, errMsgBase+task.Fullname()+"'s until property is before launching phase")
					} else {
						entity.CloseAt = until
					}
				}

				return nil

			}); err != nil {
				return errors.WithMessage(err, errMsgBase+task.Fullname())
			}
		}
	}

	// TODO: BEGIN NEXT

	// relation setting
	for iPhase := range r.Phases {
		phase := r.Phases[iPhase]

		for iTask := range phase.Tasks {
			task := &phase.Tasks[iTask]

			task.Task.Element[task_element.StdOut].Sendto = make([]entities.Relation, 0)
			task.Task.Element[task_element.StdErr].Sendto = make([]entities.Relation, 0)

			// standard output sender task's sendto property
			if task.StdOut != nil {

				task.Element[task_element.StdOut] = *task.StdOut.Element

				// when stdout property is existed
				for iRel := range task.StdOut.SendTo {
					rel := task.StdOut.SendTo[iRel]

					// validate task's name
					if phaseName, taskName, taskElem, err := entities.ParseTaskElem(rel, task.Phasename, task_element.StdIn); err != nil {
						return errors.WithMessage(err, errMsgBase+task.Fullname()+"'s sendto property: "+rel)

						// check reciever task is exist or not
					} else if reciever, err := r.Rehearsal.Task(entities.TaskFullname(phaseName, taskName)); err != nil {
						return errors.WithMessage(err, errMsgBase+task.Fullname()+"'s sendto property: "+rel)
					} else {

						// check reciever task's element
						switch taskElem {
						case task_element.StdIn:
							task.Task.AddRelation(task_element.StdOut, reciever, taskElem)
						default:
							errors.New(errMsgBase + task.Fullname() + "' sendto property: sendto task's element is invalid as reciever one")
						}
					}
				}
			}

			// standerd error output sender task's sendto property
			if task.StdErr != nil {

				task.Element[task_element.StdErr] = *task.StdErr.Element

				// when stderr property is existed
				for iRel := range task.StdErr.SendTo {
					rel := task.StdErr.SendTo[iRel]

					// validate task's name
					if phaseName, taskName, taskElem, err := entities.ParseTaskElem(rel, task.Phasename, task_element.StdIn); err != nil {
						return errors.WithMessage(err, errMsgBase+task.Fullname()+"'s sendto property: "+rel)

						// check reciever task is exist or not
					} else if reciever, err := r.Rehearsal.Task(entities.TaskFullname(phaseName, taskName)); err != nil {
						return errors.WithMessage(err, errMsgBase+task.Fullname()+"'s sendto property: "+rel)
					} else {

						// check reciever task's element
						switch taskElem {
						case task_element.StdIn:
							task.Task.AddRelation(task_element.StdErr, reciever, taskElem)
						default:
							errors.New(errMsgBase + task.Fullname() + "' sendto property: sendto task's element is invalid as reciever one")
						}
					}
				}
			}
		}
	}
	return nil
}
