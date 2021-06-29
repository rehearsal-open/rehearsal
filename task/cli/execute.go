package cli

import (
	"bytes"
	"log"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/packet/stdout"
)

func (t *Task) RunInit() error {

	t.progress = make(chan error)
	t.waiter = make(chan error)
	t.killer = make(chan error)
	go t.execute()

	return errors.WithStack(<-t.waiter)
}

func (t *Task) RunWait() error {
	t.progress <- nil
	return <-t.waiter
}

func (t *Task) Kill() {

}

func (t *Task) Finalize() {
	close(t.in)
	close(t.killer)
}

func (t *Task) execute() {

	// intitialize begin

	// definition variable
	bufout := bytes.NewBuffer(make([]byte, 0))
	buferr := bytes.NewBuffer(make([]byte, 0))
	exitOut := make(chan error)
	exitIn := make(chan error)
	exitErr := make(chan error)
	exitWait := make(chan error)

	log.Println(t.cmd)

	outputKilled := false
	var outputKilledLock sync.Mutex

	stdin, _ := t.cmd.StdinPipe()
	t.cmd.Stdout = bufout
	t.cmd.Stderr = buferr

	WaitForOutputKill := func() {
		outputKilledLock.Lock()
		defer outputKilledLock.Unlock()

		if !outputKilled {
			outputKilled = true
			exitOut <- nil
			exitErr <- nil
			for {
				time.Sleep(10 * time.Millisecond)
				if err, exist := <-exitOut; !exist {
					break
				} else {
					exitOut <- err
				}
			}
			for {
				time.Sleep(10 * time.Millisecond)
				if err, exist := <-exitErr; !exist {
					break
				} else {
					exitErr <- err
				}
			}
		}
	}

	t.waiter <- nil
	for isContinue := true; isContinue; {
		select {
		case <-t.progress:
			isContinue = false
		case <-t.killer:
			WaitForOutputKill()
			t.killed = true
			return
		default:

			time.Sleep(time.Duration(t.engine.Config().SyncMs))

		}
	}

	// run
	if err := t.cmd.Start(); err != nil {
		WaitForOutputKill()
		t.waiter <- errors.WithStack(err)
	}

	// input listener
	go func() {
		isContinue := true
		for {
			select {
			case <-exitIn:
				isContinue = false
			case input := <-t.in:
				if t != nil {
					if bytes, err := []byte(input.GetString()), error(nil); err != nil {
					} else {
						stdin.Write(bytes)
					}
				}
			default:
				if !isContinue {
					defer close(exitIn)
					return
				} else {
					time.Sleep(time.Duration(t.taskConf.SyncMs))
				}
			}
		}
	}()

	// output listener
	go func() {
		isContinue := true
		preLen := 0
		for {
			select {
			case <-exitOut:
				isContinue = false
			default:
				crtLen := bufout.Len()

				if skipped := true; crtLen > preLen {
					if bytes := bufout.Bytes()[preLen:crtLen]; bytes[len(bytes)-1] != 0 {
						skipped = false

						if str, err := t.toStr(bytes); err != nil {
						} else {
							t.sendOut(str)
						}
						preLen = crtLen
					}

				} else if skipped && !isContinue {
					defer close(exitOut)
					return
				} else {
					time.Sleep(time.Duration(t.taskConf.SyncMs))
				}
			}
		}
	}()

	// error output listener
	go func() {
		isContinue := true
		preLen := 0
		for {
			select {
			case <-exitErr:
				isContinue = false
			default:
				crtLen := buferr.Len()

				if skipped := true; crtLen > preLen {

					if bytes := buferr.Bytes()[preLen:crtLen]; bytes[len(bytes)-1] != 0 {
						skipped = false

						if str, err := t.toStr(bytes); err != nil {
						} else {
							t.sendErr(str)
						}
						preLen = crtLen
					}

				} else if skipped && !isContinue {
					defer close(exitErr)
					return
				} else {
					time.Sleep(time.Duration(t.taskConf.SyncMs))
				}
			}
		}
	}()

	// wait listener
	go func() {
		exitWait <- t.cmd.Wait()
	}()

	isContinue := []bool{true, true}
	for isContinue[1] {
		select {
		case err := <-exitWait:
			isContinue[0] = false
			WaitForOutputKill()
			t.waiter <- errors.WithStack(err)

		case <-t.killer:
			if isContinue[0] {
				t.cmd.Process.Kill()
			}
			t.killed = true
		default:
			if !isContinue[0] {
				isContinue[1] = false
			} else {
				time.Sleep(time.Duration(t.taskConf.SyncMs))
			}
		}
	}
	return
}

func (t *Task) sendOut(data string) {

	for _, ch := range t.out {
		ch <- stdout.Packet{
			Name: t.taskConf.Name,
			Data: data,
		}
	}
}

func (t *Task) sendErr(data string) {
	for _, ch := range t.err {
		ch <- stdout.Packet{
			Name: t.taskConf.Name,
			Data: data,
		}
	}
}
