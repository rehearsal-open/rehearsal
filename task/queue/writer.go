// task/queue/writer.go
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

import "github.com/rehearsal-open/rehearsal/entities"

// Make Writer instance.
func MakeWriter(writeTo *Reader) Writer {
	writeTo.lock.Lock()
	defer writeTo.lock.Unlock()

	// increment registered writer's count
	writeTo.nWriter++

	return &__Writer{
		parent: writeTo,
	}
}

// Write data, in this function cloning bytes array.
func (writer *__Writer) Write(elem *entities.Element, bytes []byte, onFinal func()) {
	defer onFinal()
	cache := make([]byte, len(bytes))
	if copied := copy(cache, bytes); copied != len(bytes) {
		panic("unfully clone")
	}
	writer.parent.__append(&__Packet{
		element: elem,
		bytes:   cache,
	})
}

func (writer *__Writer) Close() {
	writer.parent.lock.Lock()
	defer writer.parent.lock.Unlock()

	// decrement registered writer's count
	writer.parent.nWriter--
}
