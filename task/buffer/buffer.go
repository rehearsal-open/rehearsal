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
)

// Make byte buffer manages in task.
func MakeBuffer(entity *entities.Element) *Buffer {
	return &Buffer{
		mutex:    &sync.Mutex{},
		task:     entity,
		packets:  make([]*packetBase, 128),
		reciever: make([]SendToRecieverBased, 0),
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
						Closed:     false,
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
	for i := range b.reciever {
		b.reciever[i].SendPacket(Packet{
			Closed:     true,
			packetBase: nil,
			offset:     0,
		})
	}
	b.running = false
}

func (b *Buffer) AppendReciever(r SendToRecieverBased) {
	b.reciever = append(b.reciever, r)
}
