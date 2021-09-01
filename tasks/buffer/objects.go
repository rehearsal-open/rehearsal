// tasks/buffer/objects.go
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

type (
	// Contains written buffer
	Buffer struct {
		mutex    *sync.RWMutex
		task     *entities.Task
		element  entities.TaskElement
		packets  []packetBase
		reciever []Reciever
		ch       chan []byte
		running  bool
	}

	// Contains sending packet
	packetBase struct {
		buffer  *Buffer
		bytes   []byte
		nClosed int
		nSend   int
	}

	// Contains packet
	Packet struct {
		*packetBase
		offset int
	}

	// Recieves Packet
	Reciever interface {
		SendPacket(p Packet)
	}
)
