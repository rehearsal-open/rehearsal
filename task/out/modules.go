// task/out/modules.go
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

package out

import (
	"time"

	"github.com/rehearsal-open/rehearsal/engine"
	"github.com/rehearsal-open/rehearsal/entity"
	"github.com/rehearsal-open/rehearsal/logger"
	"github.com/rehearsal-open/rehearsal/packet/stdout"
)

type Task struct {
	engine      engine.RehearsalEngine
	config      *entity.TaskConfig
	logger      *logger.Logger
	in          chan stdout.Packet
	exitRoutine chan error
	killed      bool
}

func (t *Task) AssignEngine(e engine.RehearsalEngine, conf *entity.TaskConfig, name string) error {
	t.engine = e
	t.config = conf
	t.config.SyncMs = e.Config().SyncMs
	t.in = make(chan stdout.Packet)
	t.killed = false

	return nil
}

func (t *Task) AssignLogger(l *logger.Logger) error {
	t.logger = l
	return nil
}

func (t *Task) Config() *entity.TaskConfig {
	return t.config
}

func (t *Task) BytesFromString(src string, sendFrom string) ([]byte, error) {
	return []byte(src), nil
}

func (t *Task) In() chan stdout.Packet {
	return t.in
}

func (t *Task) RunInit() error {

	t.exitRoutine = make(chan error)
	go t.routine()
	return nil
}

func (t *Task) RunWait() error {
	return nil
}

func (t *Task) Kill() {
	t.exitRoutine <- nil
	t.killed = true
}

func (t *Task) Finalize() {

	if !t.killed {
		t.exitRoutine <- nil
	}

	for {
		time.Sleep(50 * time.Millisecond)
		if err, exist := <-t.exitRoutine; !exist {
			close(t.in)
			return
		} else {
			t.exitRoutine <- err
		}
	}
}

func (t *Task) routine() {
	isContinue := true
	for {
		select {
		case packet, exist := <-t.in:
			if exist {
				t.logger.PacketPrint(&packet)
			} else {
				t.logger.SystemPrint("packet is channel is closed")
			}
		case <-t.exitRoutine:
			isContinue = false
		default:
			if !isContinue {
				defer close(t.exitRoutine)
				return
			} else {
				time.Sleep(time.Duration(t.config.SyncMs))
			}
		}
	}
}
