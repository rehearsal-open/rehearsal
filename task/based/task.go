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

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_state"
	"github.com/rehearsal-open/rehearsal/task"
	"github.com/rehearsal-open/rehearsal/task/queue"
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

	// initialize task's element
	for i, l := 0, task_element.Len; i < l; i++ {

		elem := task_element.Enum(i)

		if impl.IsSupporting(elem) {
			switch elem {
			case task_element.StdIn:
				basis.inputs[i] = &inputElem{
					taskElement: &taskElement{
						internalTask: basis,
						element:      &entity.Element[i],
					},
					queue: queue.MakeReader(),
				}
			case task_element.StdOut, task_element.StdErr:
				basis.outputs[i] = &outputElem{
					taskElement: &taskElement{
						internalTask: basis,
						element:      &entity.Element[i],
					},
					writer: queue.MakeSenders(&entity.Element[i]),
				}
			}
		}
	}

	return basis
}

func (basis *internalTask) IsSupporting(elem task_element.Enum) bool {
	return basis.impl.IsSupporting(elem)
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

// Begin main task.
func (basis *internalTask) BeginTask() error {
	basis.lock.Lock()
	defer basis.lock.Unlock()

	// check task's running state
	if basis.mainstate == task_state.Closed || basis.mainstate == task_state.Finalized {
		return task.ErrAlreadyClosed // already closed
	} else if basis.mainstate == task_state.Running {
		return task.ErrAlreadyRun // already run
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
		basis.outputs[i], basis.inputs[i] = nil, nil
	}

	basis.mainstate = task_state.Finalized
}

// Append reciever to selected sender element.
func (basis *internalTask) Connect(senderElem task_element.Enum, recieverElem task_element.Enum, reciever task.Task) error {

	basis.lock.Lock()
	defer basis.lock.Unlock()

	var recieverBased *internalTask

	// check whether reciever is support based system
	if rec, ok := reciever.(Synthesized); !ok {
		panic("reciever task is unsupported")
	} else {
		recieverBased = rec.based()
	}

	// check whether element is supported or not
	if basis.outputs[senderElem] == nil {
		return task.ErrNotSupportingElement
	} else if recieverBased.inputs[recieverElem] == nil {
		return task.ErrNotSupportingElement
	}

	// check whether task has already run or not
	if basis.mainstate != task_state.Waiting {
		return task.ErrAlreadyRun
	} else {
		basis.outputs[senderElem].writer.AppendWriter(queue.MakeWriter(recieverBased.inputs[recieverElem].queue))
		return nil
	}
}

func (basis *internalTask) based() *internalTask {
	return basis
}

// Gets io.Writer using sender object.
func (basis *internalTask) Writer(elem task_element.Enum) (io.Writer, error) {
	if basis.outputs[elem] == nil {
		return nil, task.ErrNotSupportingElement
	} else {
		return basis.outputs[elem].writer, nil
	}
}

// Listen reciever element and begin to manage packet.
func (basis *internalTask) ListenStart(callback [task_element.Len]ImplCallback) {

	for i, l := 0, task_element.Len; i < l; i++ {
		if basis.inputs[i] != nil {

			// check whether callback isn't empty or not
			if callback[i] == nil {
				panic("task is supported, but callback is nil")
			}

			// begin to listen element goroutine
			go func(elem int) {
				for isContinue := true; isContinue; {
					basis.inputs[elem].queue.Read(func(e *entities.Element, b []byte) {
						if e != nil {
							callback[elem].Recieve(e, b)
						} else {
							callback[elem].OnFinal()
							isContinue = false
						}
					})
				}
			}(i)
		}
	}
}

func (basis *internalTask) Close(err error) {

	basis.lock.Lock()
	defer basis.lock.Unlock()

	for i, l := 0, task_element.Len; i < l; i++ {
		if basis.inputs[i] != nil {
			basis.inputs[i].queue.Close()
		}
		if basis.outputs[i] != nil {
			// time.Sleep(50 * time.Millisecond)
			basis.outputs[i].writer.Release()
		}
	}

	basis.mainstate = task_state.Closed
	close(basis.closed)
}
