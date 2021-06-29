package cli

import (
	"log"
	"os/exec"

	"github.com/rehearsal-open/rehearsal/engine"
	"github.com/rehearsal-open/rehearsal/entity"
	"github.com/rehearsal-open/rehearsal/packet/stdout"
	"github.com/rehearsal-open/rehearsal/task"
)

type Task struct {
	engine   engine.RehearsalEngine
	taskConf *entity.TaskConfig
	cmd      *exec.Cmd

	in       chan stdout.Packet
	out      []chan stdout.Packet
	err      []chan stdout.Packet
	toStr    func([]byte) (string, error)
	fromStr  func(string) ([]byte, error)
	killed   bool
	progress chan error
	waiter   chan error
	killer   chan error
}

func (t *Task) AssignEngine(e engine.RehearsalEngine, config *entity.TaskConfig, name string) error {
	t.engine = e
	t.taskConf = config
	t.in = make(chan stdout.Packet)
	t.out = make([]chan stdout.Packet, 0)
	t.err = make([]chan stdout.Packet, 0)
	t.toStr = func(src []byte) (string, error) { return string(src), nil }
	t.fromStr = func(src string) ([]byte, error) { return []byte(src), nil }
	t.killed = false

	t.cmd = exec.Command(t.taskConf.ExecPath, t.taskConf.Args...)
	t.cmd.Dir = e.Config().Dir
	log.Println(t.cmd.Dir)
	return nil
}

func (t *Task) Config() *entity.TaskConfig { return t.taskConf }

func (t *Task) AppendTaskAsOut(sendTo task.RecieverTask) error {
	t.out = append(t.out, sendTo.In())
	return nil
}

func (t *Task) AppendTaskAsErr(sendTo task.RecieverTask) error {
	t.err = append(t.err, sendTo.In())
	return nil
}

func (t *Task) In() chan stdout.Packet {
	return t.in
}

func (t *Task) BytesFromString(src string, sendFrom string) ([]byte, error) {
	return []byte(src), nil
}

func (t *Task) BytesToString(src []byte, sendFrom string) (string, error) {
	return string(src), nil
}
