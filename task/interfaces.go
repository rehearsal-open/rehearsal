package task

import (
	"github.com/rehearsal-open/rehearsal/entity"
	. "github.com/rehearsal-open/rehearsal/packet/task"
)

type Task interface {
	AssignConfig(config entity.Conf, name string) error
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
