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
		&Phase{
			Name:  "__init",
			Tasks: []Task{},
		},
	}
	finalPhases := []*Phase{
		&Phase{
			Name:  "__finl",
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

				return nil

			}); err != nil {
				return errors.WithMessage(err, errMsgBase+task.Fullname())
			}
		}
	}

	// TODO: BEGIN NEXT

	// relation setting
	for iPhase := range r.Phases {
		phase := &r.Phases[iPhase]

		for iTask := range phase.Tasks {
			task := &phase.Tasks[iTask]

			for iRel := range task.SendTo {
				rel := task.SendTo[iRel]

				if err := shortNameRegexp(rel); err == nil {
					rel = task.Phasename + "::" + rel
				} else if err := fullNameRegexp(rel); err != nil {
					return err
				}

				if sendto, err := r.Rehearsal.Task(rel); err != nil {
					return err
				} else {
					task.Task.AddRelation(entities.Relation{
						Reciever:        sendto,
						ElementSender:   task_element.StdOut,
						ElementReciever: task_element.StdIn,
					})
				}
			}
		}
	}
	return nil
}
