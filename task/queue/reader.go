// task/queue/reader.go
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

package queue

import (
	"sync"
	"time"

	"github.com/rehearsal-open/rehearsal/entities"
)

// Make Reader instance.
func MakeReader() *Reader {
	return &Reader{
		pool:          make([]*__Packet, 1024),
		nContain:      0,
		nWriter:       0,
		readPacketPos: 0,
		lock:          sync.Mutex{},
		onRecieve:     make(chan __Packet),
		isWaiting:     false,
		isClosed:      false,
	}
}

// Read next data set, if entities.
// Element argument in callback is nil, read function will not append bytes.
func (reader *Reader) Read(callback func(*entities.Element, []byte)) {

	// when totally closed
	if reader.isClosed || reader.onRecieve == nil {
		callback(nil, []byte{})
	}

	// get cached buffers
	packet, C := reader.__read()
	if packet == nil {

		if reader.nWriter > 0 {

			// wait access
			packet, exist := <-C
			if !exist {
				// on reader is closed, empty bytes.
				reader.lock.Lock()
				defer reader.lock.Unlock()

				// finalize
				reader.isWaiting = false
				reader.isClosed = true

				// call ended
				callback(nil, []byte{})
				return
			} else {
				callback(packet.element, packet.bytes)
				return
			}
		} else {
			// when registered running writer is nothing

			// call ended
			callback(nil, []byte{})
			return
		}

	}

	func() {
		reader.lock.Lock()
		defer reader.lock.Unlock()

		callback(packet.element, packet.bytes)

		// increment read index
		reader.readPacketPos++
		reader.readPacketPos %= len(reader.pool)
		reader.nContain--
		if reader.nContain < 0 {
			panic("reader.nContain is lower than 0")
		}

	}()

}

func (reader *Reader) __read() (*__Packet, chan __Packet) {

	reader.lock.Lock()
	defer reader.lock.Unlock()

	if reader.nContain > 0 {
		// successfully
		return reader.pool[reader.readPacketPos%len(reader.pool)], nil
	} else {
		// this doesn't have unread packets,
		// make isWaiting true.
		reader.isWaiting = true
		return nil, reader.onRecieve
	}
}

func (reader *Reader) Close() {

	for reader.nContain > 0 {
		time.Sleep(time.Millisecond)
	}

	reader.lock.Lock()
	defer reader.lock.Unlock()

	// close
	close(reader.onRecieve)
	reader.isClosed = true
}

func (reader *Reader) __append(packet *__Packet) {

	reader.lock.Lock()
	defer reader.lock.Unlock()

	if reader.isClosed {
		return
	} else {
		if reader.nContain >= len(reader.pool) {
			// no capacity to append packet

			buffer := make([]*__Packet, len(reader.pool)*2)
			if prehalf := copy(buffer, reader.pool[reader.readPacketPos:]); prehalf != len(reader.pool)-reader.readPacketPos {
				panic("copied size is unexpected")
			} else if sufhalf := copy(buffer[prehalf:], reader.pool[:reader.readPacketPos]); sufhalf != reader.readPacketPos {
				panic("copied size is unexpected")
			}

			reader.pool = buffer
			reader.readPacketPos = 0
		}

		if reader.isWaiting {
			reader.onRecieve <- *packet
			reader.isWaiting = false
		} else {
			// append packet
			reader.pool[(reader.readPacketPos+reader.nContain)%len(reader.pool)] = packet

			reader.nContain++

		}

	}
}
