package task

import (
	"github.com/rehearsal-open/rehearsal/engine"
	. "github.com/rehearsal-open/rehearsal/packet/task"
)

type Task interface {
	AssignEngine(config engine.RehearsalEngine, name string) error
	AppendOutPipe(Task) error
	AppendErrPipe(Task) error
	InputChan() chan Packet
	Initialize() error
	Wait() error
	Finalize() error
	Kill()
	BytesToString([]byte) (string, error)
	BytesFromString(string) ([]byte, error)
}
