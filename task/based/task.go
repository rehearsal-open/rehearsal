// task/based/task.go
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

package based

import (
	"io"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_state"
	"github.com/rehearsal-open/rehearsal/task/buffer"
)

// Make task basis.
// This function should be called from implemented task's make function.
func MakeBasis(entity *entities.Task, impl TaskImpl) Task {

	// initialize shared member
	basis := &internalTask{
		impl:      impl,
		entity:    entity,
		mainstate: task_state.Waiting,
		closed:    make(chan error, 1),
		lock:      sync.Mutex{},
	}

	support := [task_element.Len]bool{false}

	for i, l := 0, task_element.Len; i < l; i++ {
		basis.elements[i] = taskElement{
			internalTask: basis,
			element:      task_element.Enum(i),
			state:        task_state.Waiting,
			sender:       nil,
			reciever:     nil,
		}
		support[i] = impl.IsSupporting(task_element.Enum(i))
	}

	// initialize reciever element
	if support[task_element.StdIn] {
		basis.elements[task_element.StdIn].reciever = make(chan buffer.Packet)
		basis.elements[task_element.StdIn].packets = []buffer.Packet{}
		basis.elements[task_element.StdIn].packetPos = 0
		basis.elements[task_element.StdIn].packetLock = &sync.Mutex{}

		go func() {
			for {
				packet, exist := <-basis.elements[task_element.StdIn].reciever

				if exist {

					basis.elements[task_element.StdIn].packetLock.Lock()
					basis.elements[task_element.StdIn].packets = append(basis.elements[task_element.StdIn].packets, packet)
					basis.elements[task_element.StdIn].packetLock.Unlock()

				} else {
					return
				}
			}
		}()

	}

	// initialize sender element
	if support[task_element.StdOut] {
		basis.elements[task_element.StdOut].sender = buffer.MakeBuffer(basis.entity, task_element.StdOut)
	}
	if support[task_element.StdErr] {
		basis.elements[task_element.StdErr].sender = buffer.MakeBuffer(basis.entity, task_element.StdErr)
	}

	return basis
}

// Gets task's configuration in entity
func (basis *internalTask) Entity() *entities.Task {
	return basis.entity
}

// Gets main task's running state.
func (basis *internalTask) MainState() task_state.Enum {
	basis.lock.Lock()
	defer basis.lock.Unlock()
	return basis.mainstate
}

// Gets selected element's running state.
func (basis *internalTask) ElementState(elem task_element.Enum) task_state.Enum {
	basis.lock.Lock()
	defer basis.lock.Unlock()
	return basis.elements[elem].state
}

// Begin main task.
func (basis *internalTask) BeginTask() error {
	basis.lock.Lock()
	defer basis.lock.Unlock()

	// check task's running state
	if basis.mainstate == task_state.Closed || basis.mainstate == task_state.Finalized {
		return ErrAlreadyClosed // already closed
	} else if basis.mainstate == task_state.Running {
		return ErrAlreadyRun // already run
	}

	// begin main task
	if err := basis.impl.ExecuteMain(basis); err != nil {
		basis.mainstate = task_state.Closed
		return errors.WithMessage(err, "failed to begin task")
	} else {
		basis.mainstate = task_state.Running
		return nil
	}
}

// Stop main task
func (basis *internalTask) StopTask() {
	basis.impl.StopMain()
	<-basis.closed
}

// Wait for main task closing.
func (basis *internalTask) WaitClosing() {
	for {
		_, exist := <-basis.closed
		if !exist {
			return
		}
	}
}

func (basis *internalTask) ReleaseResource() {

	// check whether main task is running or not
	if state := basis.MainState(); state == task_state.Running {
		basis.StopTask()
		basis.WaitClosing()
	}

	basis.lock.Lock()
	defer basis.lock.Unlock()

	// release resource
	for i, l := 0, task_element.Len; i < l; i++ {
		basis.elements[i].sender = nil
		basis.elements[i].reciever = nil
		basis.elements[i].state = task_state.Finalized
	}

	basis.mainstate = task_state.Finalized
}

// Append reciever to selected sender element.
func (basis *internalTask) AppendReciever(sender task_element.Enum, reciever buffer.Reciever) error {

	basis.lock.Lock()
	defer basis.lock.Unlock()

	// check whether element is supported or not
	if basis.elements[sender].sender == nil {
		return ErrNotSupportingElement
	}

	// check whether task has already run or not
	if basis.mainstate != task_state.Waiting {
		return ErrAlreadyRun
	} else {
		basis.elements[sender].sender.AppendReciever(reciever)
		return nil
	}

}

// Get reciever selected by task element.
func (basis *internalTask) Reciever(reciever task_element.Enum) (buffer.Reciever, error) {

	basis.lock.Lock()
	defer basis.lock.Unlock()

	// check whether element is supported or not
	if basis.elements[reciever].reciever == nil {
		return nil, ErrNotSupportingElement
	}

	// check whether task has already run or not
	if basis.mainstate != task_state.Waiting {
		return nil, ErrAlreadyRun
	} else {
		return &basis.elements[reciever], nil
	}
}

// Gets io.Writer using sender object.
func (basis *internalTask) Writer(elem task_element.Enum) (io.Writer, error) {
	if basis.elements[elem].sender == nil {
		return nil, ErrNotSupportingElement
	} else {
		return basis.elements[elem].sender, nil
	}
}

// Listen reciever element and begin to manage packet.
func (basis *internalTask) ListenStart(callback [task_element.Len]Reciever) {

	for i, l := 0, task_element.Len; i < l; i++ {
		if basis.elements[i].reciever != nil {

			// check whether callback isn't empty or not
			if callback[i] == nil {
				panic("task is supported, but callback is nil")
			}

			// begin to listen element goroutine
			go func(elem int) {
				for {

					for len(basis.elements[elem].packets) <= basis.elements[elem].packetPos {
						time.Sleep(time.Millisecond)
						if basis.elements[elem].packets == nil {
							return
						}
					}

					basis.elements[elem].lock.Lock()
					callback[elem](&basis.elements[elem].packets[basis.elements[elem].packetPos])
					basis.elements[elem].packetPos++
					basis.elements[elem].lock.Unlock()

				}
			}(i)
			basis.elements[i].state = task_state.Running
		}

		if basis.elements[i].sender != nil {
			basis.elements[i].sender.Begin()
			basis.elements[i].state = task_state.Running
		}
	}
}

func (basis *internalTask) Close(err error) {

	basis.lock.Lock()
	defer basis.lock.Unlock()

	for i, l := 0, task_element.Len; i < l; i++ {
		if basis.elements[i].reciever != nil {
			for len(basis.elements[i].reciever) > 0 {
			}
			close(basis.elements[i].reciever)
		}
		if basis.elements[i].sender != nil {
			basis.elements[i].sender.Close()
		}
		basis.elements[i].state = task_state.Closed
	}

	basis.mainstate = task_state.Closed
	close(basis.closed)
}
