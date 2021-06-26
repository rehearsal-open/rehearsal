package console

import (
	"context"
	"os/exec"
	"sync"

	"github.com/rehearsal-open/rehearsal/executer"
)

type Packet struct {
	sendFrom string
	data     string
	priority executer.ErrorPriority
}

type Execute struct {
	name     string
	cmd      *exec.Cmd
	cancel   context.CancelFunc
	in       chan executer.Packet
	out      [](chan executer.Packet)
	err      [](chan executer.Packet)
	killall  chan string
	isrun    bool
	mutex    sync.Mutex
	toStr    func([]byte) (string, error)
	fromStr  func(string) ([]byte, error)
	progress chan string
	waiter   chan string
}

func ConsoleExecuteMaker(name string, args ...string) *Execute {
	ctx, cancel := context.WithCancel(context.Background())
	res := Execute{
		name:   name,
		cmd:    exec.CommandContext(ctx, name, args...),
		cancel: cancel,
		in:     make(chan executer.Packet),
		out:    make([]chan executer.Packet, 0),
		err:    make([]chan executer.Packet, 0),
		isrun:  false,
		toStr: func(src []byte) (string, error) {
			return string(src), nil
		},
		fromStr: func(src string) ([]byte, error) {
			return []byte(src), nil
		},
	}
	return &res
}
