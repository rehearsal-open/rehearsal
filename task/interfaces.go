// task/interfaces.go
// Copyright (C) 2021  Kasai Koji

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

package task

import (
	"github.com/rehearsal-open/rehearsal/engine"
	"github.com/rehearsal-open/rehearsal/entity"
	"github.com/rehearsal-open/rehearsal/packet/stdout"
)

type Task interface {
	AssignEngine(engine engine.RehearsalEngine, taskConf *entity.TaskConfig, name string) error

	// get task config
	Config() *entity.TaskConfig

	// call as single thread
	RunInit() error

	// call as goroutine, run start and return after task is stop.
	RunWait() error

	// after kill, must call finalize
	Kill()

	// call when all tasks are stopped
	Finalize()
}

type RecieverTask interface {
	Task
	In() chan stdout.Packet
	BytesFromString(src string, sendFrom string) ([]byte, error)
}

type OutTask interface {
	Task
	AppendTaskAsOut(RecieverTask) error
	BytesToString(src []byte, sendTo string) (string, error)
}

type ErrTask interface {
	Task
	AppendErrAsErr(RecieverTask) error
	BytesToString(src []byte, sendTo string) (string, error)
}
