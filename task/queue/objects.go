// task/queue/objects.go
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

	"github.com/rehearsal-open/rehearsal/entities"
)

type (
	__Packet struct {
		element *entities.Element
		bytes   []byte
	}

	// Packet query reciever side instance.
	Reader struct {
		// Reader[i] = pool[((i -))]
		pool []*__Packet
		// The number of packets containing now.
		nContain int
		// The position of numRPacket's index in pools.
		readPacketPos int
		// access mutex
		lock sync.Mutex
		// Use when wait for new appending. When isWaiting is false, writer  doesn't need to use this.
		onRecieve chan __Packet
		// Whether reader is waiting or not.
		isWaiting bool
		// Whether reader is closed.
		isClosed bool
	}

	// Packet query sender side instance.
	Writer struct {
		parent *Reader
	}

	// Parallel data sender.
	// Use for convert from io.Writer() to QueryWriter().
	Senders struct {
		// Use element entity.
		elem *entities.Element
		// Parallel data writer
		sendto []__Sender
		// Parallel mutex, using in sending data and close writer.
		parallelLock sync.WaitGroup
		// Access lock, Ban to multiple call function.
		accessLock sync.Mutex
		// Cache chan
		cacheChan chan []byte
	}

	__Sender struct {
		parent *Senders
		*Writer
		conn chan []byte
	}
)
