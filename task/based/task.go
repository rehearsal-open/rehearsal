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
func MakeBasis(entity *entities.Task, impl TaskImpl) *Task {

	// initialize shared member
	basis := &Task{
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
						Task:    basis,
						element: &entity.Element[i],
					},
					queue: queue.MakeReader(),
				}
			case task_element.StdOut, task_element.StdErr:
				basis.outputs[i] = &outputElem{
					taskElement: &taskElement{
						Task:    basis,
						element: &entity.Element[i],
					},
					writer: queue.MakeSenders(&entity.Element[i]),
				}
			}
		}
	}

	return basis
}

func (basis *Task) IsSupporting(elem task_element.Enum) bool {
	return basis.impl.IsSupporting(elem)
}

// Gets task's configuration in entity
func (basis *Task) Entity() *entities.Task {
	return basis.entity
}

// Gets main task's running state.
func (basis *Task) MainState() task_state.Enum {
	basis.lock.Lock()
	defer basis.lock.Unlock()
	return basis.mainstate
}

// Begin main task.
func (basis *Task) BeginTask() error {
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
func (basis *Task) StopTask() {
	basis.impl.StopMain()
	<-basis.closed
}

// Wait for main task closing.
func (basis *Task) WaitClosing() {
	for {
		_, exist := <-basis.closed
		if !exist {
			return
		}
	}
}

func (basis *Task) ReleaseResource() {

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

func (basis *Task) GetInput(elem task_element.Enum) *queue.Reader {
	if basis.inputs[elem] == nil {
		return nil
	} else {
		return basis.inputs[elem].queue
	}
}

func (basis *Task) GetOutput(elem task_element.Enum) *queue.Senders {
	if basis.outputs[elem] == nil {
		return nil
	} else {
		return basis.outputs[elem].writer
	}
}

// Append reciever to selected sender element.
func (basis *Task) Connect(senderElem task_element.Enum, recieverElem task_element.Enum, reciever task.Task) error {

	basis.lock.Lock()
	defer basis.lock.Unlock()

	var (
		recieverBased queue.Task
	)

	// check whether reciever is support based system
	if rec, ok := reciever.(queue.Task); !ok {
		panic("reciever task is unsupported")
	} else {
		recieverBased = rec
	}

	// check whether task has already run or not
	if basis.mainstate != task_state.Waiting || reciever.MainState() != task_state.Waiting {
		return task.ErrAlreadyRun
	}

	// check whether element is supported or not
	if from := basis.GetOutput(senderElem); from == nil {
		return task.ErrNotSupportingElement
	} else if to := recieverBased.GetInput(recieverElem); to == nil {
		return task.ErrNotSupportingElement
	} else {
		from.AppendWriter(queue.MakeWriter(to))
		return nil
	}
}

// Gets io.Writer using sender object.
func (basis *Task) Writer(elem task_element.Enum) (io.Writer, error) {
	if basis.outputs[elem] == nil {
		return nil, task.ErrNotSupportingElement
	} else {
		return basis.outputs[elem].writer, nil
	}
}

func (basis *Task) Close(err error) {

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
