// task/impl/net/task.go
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

package net

import (
	"io"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/task/based"
	"github.com/rehearsal-open/rehearsal/task/wrapper/listen"
)

func (t *__task) IsSupporting(elem task_element.Enum) bool {
	return [task_element.Len]bool{
		true, true, false,
	}[elem]
}

func (t *__task) ExecuteMain(args based.MainFuncArguments) error {

	var stdOut io.Writer

	if conn, err := net.DialTimeout(t.Entity().Kind, t.Address, 10*time.Second); err != nil {
		return errors.WithStack(err)
	} else {
		t.Conn = conn
	}

	if stdout, err := args.Writer(task_element.StdOut); err != nil {
		return errors.WithStack(err)
	} else {
		stdOut = stdout
	}

	sync := time.Duration(t.Detail.SyncSec * float64(int64(time.Second)))
	closer := listen.SyncIoPipe(t.Conn, stdOut, sync, func(e error) {
		if e == net.ErrClosed {
			t.StopMain()
		}
	})

	listen.Listen(t, task_element.StdIn, t.Conn, func(e error) {
		if e == net.ErrClosed {
			t.StopMain()
		}
	}, nil)

	go func() {

		<-t.close
		closer <- nil
		t.Conn.Close()

		args.Close(nil)
	}()

	return nil
}

func (t *__task) StopMain() {
	close(t.close)
}
