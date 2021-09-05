// task/queue/senders.go
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

// Make parallel senders instance.
func MakeSenders(senderElem *entities.Element) *Senders {

	senders := Senders{
		elem:         senderElem,
		sendto:       make([]__Sender, 0),
		parallelLock: sync.WaitGroup{},
		accessLock:   sync.Mutex{},
	}

	return &senders
}

func (senders *Senders) Write(src []byte) (written int, err error) {
	senders.accessLock.Lock()
	defer senders.accessLock.Unlock()

	// get number of parallel
	numParallel := len(senders.sendto)
	senders.parallelLock.Add(numParallel)

	// send data to all parallel goroutine
	for i := range senders.sendto {
		senders.sendto[i].conn <- src
	}

	// wait for done all parallel goroutine's task
	senders.parallelLock.Wait()

	return len(src), nil

}

// Release writer and stop goroutine.
func (senders *Senders) Release() {
	senders.accessLock.Lock()
	defer senders.accessLock.Unlock()

	// send data to all parallel goroutine
	numParallel := len(senders.sendto)
	senders.parallelLock.Add(numParallel)

	for i := range senders.sendto {
		close(senders.sendto[i].conn)
	}

	senders.parallelLock.Wait()
}

func (senders *Senders) AppendWriter(writer *Writer) {
	senders.accessLock.Lock()
	defer senders.accessLock.Unlock()

	appended := len(senders.sendto)

	// append writer
	senders.sendto = append(senders.sendto, __Sender{
		parent: senders,
		Writer: writer,
		conn:   make(chan []byte),
	})

	go senders.sendto[appended].__routine()
}

func (sender *__Sender) __routine() {
	for {
		bytes, exist := <-sender.conn
		if exist {
			sender.Write(sender.parent.elem, bytes)
			sender.parent.parallelLock.Done()
		} else {
			sender.parent.parallelLock.Done()
			return
		}
	}
}
