package run

import (
	"os/exec"

	"github.com/pkg/errors"
)

// Constractor of run.Exec
func ExecMaker(id int, execPath string, args []string, dir string) Exec {

	env := make([]string, 0)

	return Exec{
		id: id,
		cmd: exec.Cmd{
			Path: execPath,
			Args: args,
			Env:  env,
			Dir:  dir,
		},
		timeoutMs: -1,
		sendTo:    make([]*(chan IOExpression), 1),
		Recieve:   make(chan IOExpression),
		errSendTo: make([]*(chan IOExpression), 1),
		state:     ExecRunOnWaiting,
	}
}

// Return timeout time of exec running (millisecond).
func (exec *Exec) GetTimeoutMs() int64 { return exec.timeoutMs }

// Set timeout time of exec running (millisecond).
// Can use only exec states are on waiting, if other states, errored.
func (exec *Exec) SetTimeoutMs(timeoutMs int64) error {

	if err := exec.RunAt(ExecRunOnWaiting, func() error {
		exec.timeoutMs = timeoutMs
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (exec *Exec) AppendInputReciever(reciever *(chan IOExpression)) error {
	if err := exec.RunAt(ExecRunOnWaiting, func() error {
		exec.sendTo = append(exec.sendTo, reciever)
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (exec *Exec) AppendErrorReciever(reciever *(chan IOExpression)) error {
	if err := exec.RunAt(ExecRunOnWaiting, func() error {
		exec.errSendTo = append(exec.errSendTo, reciever)
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
