package run

import (
	"github.com/pkg/errors"
)

// Judge the state of exec running, if flag is true, call the action(flag is false, create the error)
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

// Start exec running under goroutine.
func (exec *Exec) Start() error {
	if err := exec.RunAt(ExecRunOnWaiting, func() error {
		go exec.start()
		return nil
	}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (exec *Exec) start() {

}

// Force stop exec running.
func (exec *Exec) Kill() error {
	if err := exec.RunAt(ExecRunOnRunning, func() error {
		exec.state = ExecRunOnKilled
		return exec.cmd.Process.Kill()
	}); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
