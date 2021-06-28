package consoleout

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/rehearsal-open/rehearsal/engine"
	"github.com/rehearsal-open/rehearsal/entity"
	"github.com/rehearsal-open/rehearsal/logger"
	. "github.com/rehearsal-open/rehearsal/packet"
	"github.com/rehearsal-open/rehearsal/task"
)

type Task struct {
	// dependancy
	engine.RehearsalEngine
	*entity.TaskConf
	*logger.PacketLogger

	// I/O chan
	in        chan Packet
	out       map[string](chan Packet)
	err       map[string](chan Packet)
	kill      chan string
	progress  chan string
	waiter    chan string
	finalized bool
	mutex     sync.Mutex
}

func (t *Task) AssignEngine(engine engine.RehearsalEngine, name string) error {
	t.RehearsalEngine = engine
	if taskConf, exist := t.Config().TasksMap[name]; !exist {
		return errors.New("cannot found registered task name in engine: " + name)
	} else {
		t.TaskConf = taskConf
	}

	return nil
}

func (t *Task) AssignPacketLogger(logger *logger.PacketLogger) error {
	t.PacketLogger = logger
	return nil
}

func (t *Task) TaskConfig() *entity.TaskConf {
	return t.TaskConf
}

func (t *Task) AppendOutPipe(reciever task.Task) error {
	recieverName := reciever.TaskConfig().Name
	if _, exist := t.out[recieverName]; exist {
		return errors.New("already exist standard output pipe: from " + t.Name + " to " + recieverName)
	}
	t.out[recieverName] = reciever.InputChan()
	return nil
}

func (t *Task) AppendErrPipe(reciever task.Task) error {
	recieverName := reciever.TaskConfig().Name
	if _, exist := t.out[recieverName]; exist {
		return errors.New("already exist error out pipe: from " + t.Name + " to " + recieverName)
	}
	t.out[recieverName] = reciever.InputChan()
	return nil
}

func (t *Task) InputChan() chan Packet {
	return t.in
}

func (t *Task) BytesFromStrig(src string) ([]byte, error) {
	return []byte(src), nil
}

func (t *Task) BytesToString(src []byte) (string, error) {
	return string(src), nil
}
