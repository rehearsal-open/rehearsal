// tasks/basis.gateways/task.go
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

package basis

import (
	"io"

	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/element"
	"github.com/rehearsal-open/rehearsal/entities/run_state"
	"github.com/rehearsal-open/rehearsal/tasks/buffer"
)

func MakeBasis(entity *entities.Task, implemented TaskImpl) Task {

	// initialize shared member
	basis := &internalTask{
		impl:       implemented,
		entity:     entity,
		numSupport: 0,
		mainstate:  run_state.Waiting,
		closed:     make(chan error),
	}

	support := [element.NumTaskElement]bool{false}

	for i, l := 0, element.NumTaskElement; i < l; i++ {
		basis.state[i] = run_state.Waiting
		basis.sender[i] = nil
		basis.reciever[i] = nil

		if basis.impl.IsSupporting(element.TaskElement(i)) {

			// when supported element
			basis.numSupport++
			support[i] = true

		}
	}

	// initialize reciever element
	if support[element.StdIn] {
		basis.reciever[element.StdIn] = make(chan buffer.Packet)
	}

	// initialize sender element
	if support[element.StdOut] {
		basis.sender[element.StdOut] = buffer.MakeBuffer(basis.entity, element.StdOut)
	}
	if support[element.StdErr] {
		basis.sender[element.StdErr] = buffer.MakeBuffer(basis.entity, element.StdErr)
	}

	return basis
}

func (basis *internalTask) MainState() run_state.RunningState {
	return basis.mainstate
}

func (basis *internalTask) ElementState(elem element.TaskElement) run_state.RunningState {
	return basis.state[elem]
}

func (basis *internalTask) BeginTask() error {
	if err := basis.impl.ExecuteMain(basis); err != nil {
		return err
	} else {
		basis.mainstate = run_state.Running
		return nil
	}
}

func (basis *internalTask) StopTask() {
	basis.impl.StopMain()
}

func (basis *internalTask) WaitClosing() {
	<-basis.closed
}

func (basis *internalTask) ReleaseResource() {

	if basis.mainstate == run_state.Running {
		basis.WaitClosing()
	}

	for i, l := 0, element.NumTaskElement; i < l; i++ {
		basis.sender[i], basis.reciever[i], basis.state[i] = nil, nil, run_state.Finalized
	}

	basis.mainstate = run_state.Finalized
}

func (basis *internalTask) AppendReciever(sender element.TaskElement, reciever buffer.Reciever) error {

	if basis.sender[sender] != nil {
		basis.sender[sender].AppendReciever(reciever)
		return nil
	} else {
		return ErrNotSupportingElement
	}
}

func (basis *internalTask) Reciever(element element.TaskElement) (buffer.Reciever, error) {
	if basis.reciever[element] != nil {
		return &Element{basis, element}, nil
	} else {
		return nil, ErrNotSupportingElement
	}
}

func (basis *internalTask) Writer(element element.TaskElement) (io.Writer, error) {
	if basis.sender[element] == nil {
		return nil, ErrNotSupportingElement
	} else {
		return basis.sender[element], nil
	}
}

func (basis *internalTask) ListenStart() error {

	// begin reciever element
	if basis.reciever[element.StdIn] != nil {
		if callback, err := basis.impl.RecieverCallback(element.StdIn); err != nil {
			return err
		} else {
			go func() {
				for {
					packet, exist := <-basis.reciever[element.StdIn]
					if exist {
						callback(packet)
					} else {
						return
					}
				}
			}()
			basis.state[element.StdIn] = run_state.Running
		}
	}

	// begin sender element
	if basis.sender[element.StdOut] != nil {
		basis.sender[element.StdOut].Begin()
		basis.state[element.StdOut] = run_state.Running
	}

	if basis.sender[element.StdErr] != nil {
		basis.sender[element.StdErr].Begin()
		basis.state[element.StdErr] = run_state.Running
	}

	return nil
}

func (basis *internalTask) Close(err error) {

	// close reciever element
	if basis.reciever[element.StdIn] != nil {
		close(basis.reciever[element.StdIn])
		basis.state[element.StdIn] = run_state.Closed
	}

	// close sender element
	if basis.sender[element.StdOut] != nil {
		basis.sender[element.StdOut].Close()
		basis.state[element.StdOut] = run_state.Closed
	}

	if basis.sender[element.StdErr] != nil {
		basis.sender[element.StdErr].Close()
		basis.state[element.StdErr] = run_state.Closed
	}

	basis.mainstate = run_state.Closed
	close(basis.closed)

}
