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
	"errors"
	"time"

	"github.com/streamwest-1629/textfilter"
)

// Make new Rehearsal instance.
func NewRehearsal() *Rehearsal {
	return &Rehearsal{
		phaseFilter: textfilter.Multiple{
			textfilter.Identifier(),
			IsDefinedName,
		},
		SyncInterval: time.Microsecond,
		phases:       make([]*Phase, 0),
		nameList:     make(map[string]int),
	}
}

// Make new Phase instance and append it to this instance.
func (r *Rehearsal) AppendPhase(name string) (phase *Phase, err error) {

	// validdate phase's name
	if _, err = r.phaseFilter.Add(
		name,
		func(passed string) (done bool, err error) {

			// when successfully filtered

			// make phase
			phase = &Phase{
				name: passed,
				taskFilter: textfilter.Multiple{
					textfilter.Identifier(),
					IsDefinedName,
				},
				SyncInterval: r.SyncInterval,
				tasks:        make([]*Task, 0),
				nameList:     make(map[string]int),
			}

			// append phase
			idx := len(r.phases)
			r.phases = append(r.phases, phase)
			r.nameList[name] = idx

			return true, nil

		}); err != nil {
		return nil, err
	}
	return
}

// Return phase object selected by order.
func (r *Rehearsal) PhaseAt(idx int) (*Phase, error) {
	if idx > -1 && idx < len(r.phases) {
		return r.phases[idx], nil
	} else {
		return nil, errors.New("phase's order is overed")
	}
}

// Return number of including phases.
func (r *Rehearsal) PhaseLen() int {
	return len(r.phases)
}

// Return phase object selected by name.
func (r *Rehearsal) Phase(name string) (*Phase, error) {
	if idx, ok := r.nameList[name]; !ok {
		return nil, errors.New("cannot found phase's name")
	} else {
		return r.phases[idx], nil
	}
}

// Make new Task instance and append it to this instance.
func (r *Rehearsal) AppendTask(phaseName string, taskName string, kind string, detail TaskDetail) (*Task, error) {
	if phase, err := r.Phase(phaseName); err != nil {
		return nil, err
	} else if task, err := phase.appendTask(taskName, kind, detail); err != nil {
		return nil, err
	} else {
		return task, nil
	}
}

// Return task object selected by name.
func (r *Rehearsal) Task(phaseName string, taskName string) (*Task, error) {
	if p, e := r.Phase(phaseName); e != nil {
		return nil, e
	} else {
		return p.Task(phaseName)
	}
}
