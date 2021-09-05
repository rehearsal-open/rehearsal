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
	"bytes"
	"io"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_state"
	"github.com/rehearsal-open/rehearsal/task"
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
					reciever:    nil,
					numSendFrom: 0,
				}
			case task_element.StdOut, task_element.StdErr:
				basis.outputs[i] = &outputElem{
					taskElement: &taskElement{
						internalTask: basis,
						element:      &entity.Element[i],
					},
					sender: nil,
				}
			}
		}
	}

	// initialize reciever element
	if basis.inputs[task_element.StdIn] != nil {
		basis.inputs[task_element.StdIn].reciever = make(chan buffer.Packet)
		basis.inputs[task_element.StdIn].packets = []buffer.Packet{}
		basis.inputs[task_element.StdIn].packetPos = [...]int{0, 0}
		basis.inputs[task_element.StdIn].packetLock = &sync.Mutex{}

		go func() {
			for {
				packet, exist := <-basis.inputs[task_element.StdIn].reciever

				if exist {

					basis.inputs[task_element.StdIn].packetLock.Lock()
					if l := len(basis.inputs[task_element.StdIn].packets); basis.inputs[task_element.StdIn].packetPos[1] >= l {
						if (l-basis.inputs[task_element.StdIn].packetPos[0])*2 < l {
							copied := copy(basis.inputs[task_element.StdIn].packets, basis.inputs[task_element.StdIn].packets[basis.inputs[task_element.StdIn].packetPos[0]:])
							basis.inputs[task_element.StdIn].packetPos = [...]int{0, copied}
						} else {
							buf := make([]buffer.Packet, (l+1)*2)
							copied := copy(buf, basis.inputs[task_element.StdIn].packets[basis.inputs[task_element.StdIn].packetPos[0]:])
							basis.inputs[task_element.StdIn].packetPos = [...]int{0, copied}
							basis.inputs[task_element.StdIn].packets = buf
						}
					}

					basis.inputs[task_element.StdIn].packets[basis.inputs[task_element.StdIn].packetPos[1]] = packet
					basis.inputs[task_element.StdIn].packetPos[1]++
					basis.inputs[task_element.StdIn].packetLock.Unlock()

				} else {
					return
				}
			}
		}()

	}

	// initialize sender element
	if basis.outputs[task_element.StdOut] != nil {
		basis.outputs[task_element.StdOut].sender = buffer.MakeBuffer(&basis.entity.Element[task_element.StdOut])
	}
	if basis.outputs[task_element.StdErr] != nil {
		basis.outputs[task_element.StdErr].sender = buffer.MakeBuffer(&basis.entity.Element[task_element.StdOut])
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
	if rec, ok := reciever.(frontTask); !ok {
		panic("reciever task is unsupported")
	} else {
		recieverBased = rec.based()
	}

	// check whether element is supported or not
	if basis.outputs[senderElem] == nil {
		return ErrNotSupportingElement
	} else if recieverBased.inputs[recieverElem] == nil {
		return ErrNotSupportingElement
	}

	// check whether task has already run or not
	if basis.mainstate != task_state.Waiting {
		return ErrAlreadyRun
	} else {
		basis.outputs[senderElem].sender.AppendReciever(recieverBased.inputs[recieverElem])
		recieverBased.inputs[recieverElem].Registered()
		return nil
	}
}

func (basis *internalTask) based() *internalTask {
	return basis
}

// Gets io.Writer using sender object.
func (basis *internalTask) Writer(elem task_element.Enum) (io.Writer, error) {
	if basis.outputs[elem] == nil {
		return nil, ErrNotSupportingElement
	} else {
		return basis.outputs[elem].sender, nil
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
				iClosed := 0
				for {

					// packet, exist := <-basis.elements[elem].reciever
					// if exist {
					// 	if packet.Closed {
					// 		if iClosed++; iClosed >= basis.elements[elem].numSendFrom {
					// 			callback[elem].OnFinal()
					// 		}
					// 	} else {
					// 		callback[elem].Recieve(&packet)
					// 	}
					// } else {
					// 	return
					// }

					for at, total := basis.inputs[elem].packetPos[0], basis.inputs[elem].packetPos[1]; total <= at; at, total = basis.inputs[elem].packetPos[0], basis.inputs[elem].packetPos[1] {
						time.Sleep(time.Millisecond)

						if basis.inputs[elem].packets == nil {
							return
						}
					}

					func() {
						basis.inputs[elem].lock.Lock()
						defer func() {
							if basis.inputs[elem].packetPos[1] > basis.inputs[elem].packetPos[0] {
								basis.inputs[elem].packetPos[0]++
							}
							basis.inputs[elem].lock.Unlock()
						}()

						if basis.inputs[elem].packetPos[1] > basis.inputs[elem].packetPos[0] {
							packet := &basis.inputs[elem].packets[basis.inputs[elem].packetPos[0]]

							if packet.Closed {
								if iClosed++; iClosed >= basis.inputs[elem].numSendFrom { // TODO: number of registered sender task
									callback[elem].OnFinal()
								}

							} else {
								bytes := bytes.NewBuffer([]byte{})
								io.Copy(bytes, packet)
								callback[elem].Recieve(packet.Sender(), bytes.Bytes())
							}
						}

					}()

				}
			}(i)
		}

		if basis.outputs[i] != nil {
			basis.outputs[i].sender.Begin()
		}
	}
}

func (basis *internalTask) Close(err error) {

	basis.lock.Lock()
	defer basis.lock.Unlock()

	for i, l := 0, task_element.Len; i < l; i++ {
		if basis.inputs[i] != nil {

			// wait for finally read
			for isContinue := true; isContinue; {
				func(i int) {
					basis.inputs[i].lock.Lock()
					defer basis.inputs[i].lock.Unlock()
					if basis.inputs[i].packetPos[1] >= basis.inputs[i].packetPos[0] {
						isContinue = false
					}

				}(i)
				if isContinue {
					time.Sleep(time.Millisecond)
				}
			}
			close(basis.inputs[i].reciever)
			basis.inputs[i].reciever = nil
			basis.inputs[i].packets = nil
		}
		if basis.outputs[i] != nil {
			basis.outputs[i].sender.Close()
		}
	}

	basis.mainstate = task_state.Closed
	close(basis.closed)
}
