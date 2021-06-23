package run

import (
	"io"

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

// Start exec running (Internal function).
func (exec *Exec) start() {
	stdin, err := exec.cmd.StdinPipe()
	if err != nil {
		// todo: error send
	}
	stdout, err := exec.cmd.StdoutPipe()
	if err != nil {
		// todo: error send
	}
	stderr, err := exec.cmd.StderrPipe()
	if err != nil {
		// todo: error send
	}
	if err := exec.cmd.Start(); err != nil {
		// todo: error send
	}

	go func() {
		for exec.cmd.ProcessState == nil {
		}
		for !exec.cmd.ProcessState.Exited() {
			// todo: read pipe loop
			buffer, err := io.ReadAll(stderr)
			if err != nil {
				// todo: error send
				//

			} else if len(buffer) > 0 {
				// todo: error send
			}

			buffer, err = io.ReadAll(stdout)
			if err != nil {
				// todo: error send
			} else if len(buffer) > 0 {
				// todo: output send to sendTo
			}

			for recieve := range exec.Recieve {
				io.WriteString(stdin, recieve.data)
			}
		}
	}()

	if err := exec.cmd.Wait(); err != nil {
		// todo: error send
		
	} else if exec.cmd.ProcessState.ExitCode == -1 {
		// todo: error send
	}

}

func (exec *Exec) sendErrorChan(err error, priority IOErrorPriority) {
	// todo: error send

	buffer := IOExpression{
		fromId: exec.id,
		data: err.Error()
		err: IOErrorProps {
			err: err,
			priority: priority,
		}
	}

	for ch := range exec.errSendTo {
		ch <- buffer
	}
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
