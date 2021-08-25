// task/buffer/buffer.go
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

package buffer

import (
	"sync"

	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
)

// Make byte buffer manages in task.
func MakeBuffer(entity *entities.Task, element task_element.Enum) *Buffer {
	return &Buffer{
		mutex:    &sync.Mutex{},
		task:     entity,
		element:  element,
		packets:  make([]*packetBase, 128),
		reciever: make([]Reciever, 0),
		ch:       make(chan []byte),
		running:  false,
	}
}

func (b *Buffer) Write(bytes []byte) (written int, err error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	clone := make([]byte, 0, len(bytes))
	clone = append(clone, bytes...)
	written = len(clone)

	b.ch <- clone

	return written, nil
}

func (b *Buffer) Begin() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.running = true
	b.nSend = len(b.reciever)

	go func() {
		for {
			clone, exist := <-b.ch
			if exist {
				base := packetBase{
					buffer:  b,
					bytes:   clone,
					nClosed: 0,
				}

				b.packets = append(b.packets, &base)

				for _, rec := range b.reciever {
					rec.SendPacket(Packet{
						packetBase: &base,
						offset:     0,
					})
				}
			} else {
				return
			}
		}
	}()
}

func (b *Buffer) Close() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	close(b.ch)
	b.running = false
}

func (b *Buffer) AppendReciever(r Reciever) {
	b.reciever = append(b.reciever, r)
}
