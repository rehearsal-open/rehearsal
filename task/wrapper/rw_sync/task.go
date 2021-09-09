// task/wrapper/rw_sync/task.go
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

package rw_sync

import (
	"sync"
	"time"

	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/task/based"
	"github.com/rehearsal-open/rehearsal/task/wrapper/listen"
)

func (ticked *TickedInput) IsSupporting(elem task_element.Enum) bool {
	return [task_element.Len]bool{true, false, false}[elem]
}

func (ticked *TickedInput) ExecuteMain(args based.MainFuncArguments) error {

	for i := range ticked.Tickers {
		switch elem := task_element.Enum(i); elem {
		case task_element.StdIn: // input
			{
				if ticked.Tickers[i] == nil {
					ticked.Tickers[i] = time.NewTicker(time.Millisecond)
				}

				closer := false
				lock := sync.Mutex{}
				cache := []byte{}

				// launch sender
				go func() {
					for {
						select {
						case <-ticked.Tickers[elem].C:
							func() {
								lock.Lock()
								defer lock.Unlock()
								if len(cache) > 0 {
									ticked.WriteCloser.Write(cache)
									cache = cache[:0]

								} else if closer {
									if len(cache) > 0 {
										ticked.WriteCloser.Write(cache)
									}
									ticked.WriteCloser.Close()
									ticked.Tickers[elem].Stop()
									ticked.closed <- nil
									return
								}
							}()
						}
					}
				}()

				// launch listener
				listen.ListenElemBytes(ticked, elem, func(_ *entities.Element, bytes []byte) {
					lock.Lock()
					defer lock.Unlock()
					cache = append(cache, bytes...)
				}, func() {
					closer = true
				})
			}
		}
	}

	go func() {
		<-ticked.close
		args.Close(nil)
	}()

	return nil
}

func (ticked *TickedInput) StopMain() {
	close(ticked.close)
	<-ticked.closed
}
