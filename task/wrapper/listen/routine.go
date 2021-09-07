// task/listen/routine.go
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

package listen

import (
	"bytes"
	"io"
	"time"

	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/task"
	"github.com/rehearsal-open/rehearsal/task/wrapper"
)

var (
	OnFinalDefault = func() {}
	OnErrDefault   = func(error) {}
)

func Listen(fromTask task.Task, fromElem task_element.Enum, destInput io.Writer, onErr func(error), onFinal func()) {

	// set default value
	if onErr == nil {
		onErr = OnErrDefault
	}
	if onFinal == nil {
		onFinal = OnFinalDefault
	}

	// get queue access
	from := wrapper.GetQueueAccess(fromTask)

	// get element queue
	reader := from.GetInput(fromElem)

	go func() {
		for isContinue := true; isContinue; {
			reader.Read(func(e *entities.Element, b []byte) {
				if e != nil {
					if _, err := io.Copy(destInput, bytes.NewBuffer(b)); err != nil {
						if err != io.EOF {
							onErr(err)
						}
					}
				} else {
					onFinal()
				}
			})
		}
	}()
}

func SyncIoPipe(src io.Reader, dest io.Writer, duration time.Duration, onErr func(error)) (closer chan error) {

	closer = make(chan error, 1)
	ticker := time.NewTicker(duration)

	go func() {
		for {
			select {
			case <-ticker.C:
				if _, err := io.Copy(dest, src); err != nil {
					if err != io.EOF {
						onErr(err)
					}
				}
			case <-closer:
				return
			}
		}
	}()
	return closer
}
