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

	// after kill, must call finalize
	Kill()

	// call when all tasks are stopped
	Finalize()
}

type RecieverTask interface {
	Task
	In() chan packet.Packet
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
}
