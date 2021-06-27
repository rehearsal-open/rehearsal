package cli

import (
	"context"
	"os/exec"
	"sync"

	"github.com/pkg/errors"

	"github.com/rehearsal-open/rehearsal/engine"
	"github.com/rehearsal-open/rehearsal/entity"
	. "github.com/rehearsal-open/rehearsal/packet/task"
	"github.com/rehearsal-open/rehearsal/task"
)

type Task struct {
	// dependancy
	engine.RehearsalEngine
	*entity.TaskConf
	cmd *exec.Cmd

	// I/O chan
	in        chan Packet
	out       map[string](chan Packet)
	err       map[string](chan Packet)
	kill      chan string
	progress  chan string
	waiter    chan string
	finalized bool
	cancel    context.CancelFunc
	mutex     sync.Mutex

	// byte/string converter
	toStr   func([]byte) (string, error)
	fromStr func(string) ([]byte, error)
}

func (t *Task) AssignEngine(engine engine.RehearsalEngine, name string) error {
	t.RehearsalEngine = engine
	if taskConf, exist := t.Config().TasksMap[name]; !exist {
		return errors.New("cannot found registered task name in engine: " + name)
	} else {
		t.TaskConf = taskConf
	}

	// todo: error check

	// input
	t.in = make(chan Packet)

	// command context set
	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel
	t.cmd = exec.CommandContext(ctx, t.Path, t.Args...)
	return nil
}

func (t *Task) AppendOutPipe(reciever task.Task) error {
	recieverName := reciever.TaskConfig().Name
	if _, exist := t.out[recieverName]; exist {
		return errors.New("already exist standard output pipe: from" + t.Name + " to " + recieverName)
	}

	t.out[recieverName] = reciever.InputChan()
	return nil
}

func (t *Task) AppendErrPipe(reciever task.Task) error {
	recieverName := reciever.TaskConfig().Name
	if _, exist := t.out[recieverName]; exist {
		return errors.New("already exist error output pipe: from" + t.Name + "to" + recieverName)
	}

	t.out[recieverName] = reciever.InputChan()
	return nil
}

func (t *Task) InputChan() chan Packet {
	return t.in
}

func (t *Task) BytesFromString(src string) ([]byte, error) {
	return t.fromStr(src)
}

func (t *Task) BytesToString(src []byte) (string, error) {
	return t.toStr(src)
}
