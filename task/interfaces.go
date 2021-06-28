package task

import (
	"github.com/rehearsal-open/rehearsal/engine"
	"github.com/rehearsal-open/rehearsal/entity"
	"github.com/rehearsal-open/rehearsal/packet"
)

type Task interface {
	TaskConfig() *entity.TaskConf
	AssignEngine(engine engine.RehearsalEngine, name string) error
	AppendOutPipe(reciever Task) error
	AppendErrPipe(reciever Task) error
	InputChan() chan packet.Packet
	Initialize() error
	Wait() error
	Finalize() error
	Kill()
	BytesToString([]byte) (string, error)
	BytesFromString(string) ([]byte, error)
}
