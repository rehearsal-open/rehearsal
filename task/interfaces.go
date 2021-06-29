package task

import (
	"github.com/rehearsal-open/rehearsal/engine"
	"github.com/rehearsal-open/rehearsal/packet"
)

type Task interface {
	AssignEngine(engine engine.RehearsalEngine, name string) error

	// call as single thread
	RunInit() error

	// call as goroutine, run start and return after task is stop.
	RunWait() error

	// call when all tasks are stopped
	Finalize()
}

type InTask interface {
	Task
	In() chan packet.Packet
}

type OutTask interface {
	Task
	AppendTaskAsOut(InTask) error
}

type ErrTask interface {
	Task
	AppendErrAsErr(InTask) error
}
