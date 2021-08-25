// parser/mapped/parser.go
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
	"github.com/streamwest-1629/convertobject"
	"github.com/streamwest-1629/textfilter"
)

func (p *Parser) Parse() (*entities.Rehearsal, error) {

	r := Rehearsal{}
	const errMsg = "cannot parse from map to object"

	var (
		phases = map[string]*Phase{}
	)

	if err := convertobject.DirectConvert(p.mapped, &r); err != nil {
		return nil, errors.WithMessage(err, errMsg)
	}

	shortNameRegexp := textfilter.RegexpMatches(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
	fullNameRegexp := textfilter.RegexpMatches(`^[a-zA-Z_][a-zA-Z0-9_]*::[a-zA-Z_][a-zA-Z0-9_]*$`)

	phaseFilter := textfilter.Multiple([]textfilter.Filter{
		textfilter.Identifier(),
		shortNameRegexp,
	})
	taskFilter := textfilter.Multiple{
		textfilter.Identifier(),
		fullNameRegexp,
	}

	// check phase name and set index
	for iPhase := range r.Phases {
		phase := &r.Phases[iPhase]
		if err := textfilter.RegisterFiltering(phaseFilter, phase.Name, func() error {

			if phase.Index < 1 {
				phase.Index = math.MaxInt32
			}
			return nil
		}); err != nil {
			return nil, errors.WithMessage(err, errMsg)
		}
	}

	// sort phase
	sort.Sort(phaseByIndex(r.Phases))

	// assign sorted result to phase
	for iPhase := range r.Phases {
		phase := &r.Phases[iPhase]

		phase.Index = iPhase

		// register phase by name
		phases[phase.Name] = phase
	}

	// assign task's details
	for iPhase := range r.Phases {
		phase := &r.Phases[iPhase]

		for iTask := range phase.Tasks {
			task := &phase.Tasks[iTask]

			task.Phasename = phase.Name

			if err := textfilter.RegisterFiltering(taskFilter, task.Fullname(), func() error {

				// set task detail data
				if err := p.DetailMaker.MakeDetail(task.Kind, task.Clone, task.Task); err != nil {
					return err
				}

				task.LaunchAt = iPhase

				r.Rehearsal.AddTask(task.Task)

				return nil

			}); err != nil {
				return nil, errors.WithMessage(err, errMsg)
			}
		}
	}

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
					return nil, err
				}

				if sendto, err := r.Rehearsal.Task(rel); err != nil {
					return nil, err
				} else {
					task.Task.AddRelation(entities.Reciever{
						Reciever: sendto,
						Element:  task_element.StdIn,
					})
				}
			}
		}
	}

	return r.Rehearsal, nil
}

type phaseByIndex []Phase

func (a phaseByIndex) Len() int           { return len(a) }
func (a phaseByIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a phaseByIndex) Less(i, j int) bool { return a[i].Index < a[j].Index }
