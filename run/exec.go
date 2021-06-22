package run

import (
	"os/exec"

	"github.com/pkg/errors"

	"github.com/rehearsal-open/rehearsal/algorithm"
)

type Exec struct {
	cmd       exec.Cmd
	timeoutMs int64
	outputBuf algorithm.SimpleQueue
	inputBuf  algorithm.SimpleQueue
	errBuf    algorithm.SimpleQueue
	state     ExecRunOn
}

type ExecRunOn int

const (
	ExecRunOnWaiting ExecRunOn = 1 << iota
	ExecRunOnRunning
	ExecRunOnKilled
)

func ExecMaker(execPath string, args []string, dir string) Exec {

	env := make([]string, 0)

	return Exec{
		cmd: exec.Cmd{
			Path: execPath,
			Args: args,
			Env:  env,
			Dir:  dir,
		},
		timeoutMs: -1,
		outputBuf: algorithm.SimpleQueueMaker(0),
		inputBuf:  algorithm.SimpleQueueMaker(0),
		errBuf:    algorithm.SimpleQueueMaker(0),
		state:     ExecRunOnWaiting,
	}
}

func (exec *Exec) RunAt(state ExecRunOn, action func() error) error {
	if exec.state&state != 0 {
		if err := action(); err != nil {
			return errors.WithStack(err)
		} else {
			return nil
		}
	} else {
		return errors.New("Exec current state is invalid")
	}
}

func (exec *Exec) GetTimeoutMs() int64 { return exec.timeoutMs }

func (exec *Exec) SetTimeoutMs(timeoutMs int64) error {

	if err := exec.RunAt(ExecRunOnWaiting, func() error {
		exec.timeoutMs = timeoutMs
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (exec *Exec) Kill() error {
	if err := exec.RunAt(ExecRunOnWaiting, func() error {
		exec.state = ExecRunOnKilled
		return exec.cmd.Process.Kill()
	}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
