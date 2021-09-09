// rehearsal-cli/cli/objects.go
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

package cli

import (
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/task/wrapper/elem_parallel"
	"github.com/rehearsal-open/rehearsal/task/wrapper/rw_sync"
	"github.com/rehearsal-open/rehearsal/task/wrapper/splitter"
)

const (
	ForeRed     = "\x1b[31m"
	ForeGreen   = "\x1b[32m"
	ForeYellow  = "\x1b[33m"
	ForeBlue    = "\x1b[34m"
	ForeMagenta = "\x1b[35m"
	ForeSyan    = "\x1b[36m"
	ForeReset   = "\x1b[39m"
	BackRed     = "\x1b[41m"
	BackGreen   = "\x1b[42m"
	BackYellow  = "\x1b[43m"
	BackBlue    = "\x1b[44m"
	BackMagenta = "\x1b[45m"
	BackSyan    = "\x1b[46m"
	BackReset   = "\x1b[49m"
)

type (
	__logger struct {
		lock   sync.Mutex
		isRead bool
		waiter sync.WaitGroup
	}
)

func Make(entity *entities.Rehearsal) *elem_parallel.ElemParallel {

	// make task's instance
	sync := rw_sync.Make(&entities.Task{}, &__logger{})
	sync.Tickers[task_element.StdIn] = time.NewTicker(100 * time.Millisecond)
	result := elem_parallel.Make(sync)
	return result

}

func InitLogger(entity *entities.Rehearsal, result *elem_parallel.ElemParallel) {
	colorSet := [...]string{
		ForeRed, ForeGreen, ForeYellow, ForeBlue, ForeMagenta, ForeSyan,
	}

	nameSize := 0
	idx := 0

	// get max length of logging task's names
	entity.ForeachTask(func(idx int, task *entities.Task) error {
		for i, l := 0, task_element.Len; i < l; i++ {
			switch elem := task_element.Enum(i); elem {
			case task_element.StdOut, task_element.StdErr:
				if element := task.Element[elem]; element.WriteLog {
					name := element.Fullname()
					if length := len(name); length > nameSize {
						nameSize = length
					}
				}
			}
		}
		return nil
	})

	// set max length
	entity.ForeachTask(func(_ int, task *entities.Task) error {
		for i, l := 0, task_element.Len; i < l; i++ {
			switch elem := task_element.Enum(i); elem {
			case task_element.StdOut, task_element.StdErr:
				if element := task.Element[elem]; element.WriteLog {
					name := element.Fullname()
					result.AppendElem(&element, &splitter.Splitter{
						SplitStr: []string{"\r\n", "\n", "\r"},
						Prefix:   colorSet[idx%len(colorSet)] + BackReset + name + strings.Repeat(" ", nameSize-len(name)) + " : ",
						Suffix:   ForeReset + BackReset + "\n",
					}, task_element.StdIn)
					idx++
				}
			}
		}
		return nil
	})

}

func (logger *__logger) Write(src []byte) (read int, err error) {

	logger.waiter.Add(1)
	defer logger.waiter.Done()

	logger.lock.Lock()
	defer logger.lock.Unlock()
	logger.isRead = true

	os.Stdout.Write(src)
	return len(src), nil
}

func (logger *__logger) Close() error {

	for con := true; con; {
		logger.waiter.Wait()
		con = func() bool {
			logger.lock.Lock()
			defer logger.lock.Unlock()
			if logger.isRead {
				logger.isRead = false
				return true
			} else {
				return false
			}
		}()
	}

	return nil
}
