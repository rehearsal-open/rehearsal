package cli

import (
	"bytes"
	"sync"
	"time"

	"github.com/rehearsal-open/rehearsal/packet/task"
)

const (
	Successfully string = "Successfully"
)

func (t *Task) Initialize() error {

	t.kill = make(chan string)
	t.progress = make(chan string)
	t.waiter = make(chan string)

	go t.execute()
	<-t.waiter
	return nil
}

func (t *Task) Wait() error {
	if !t.finalized {
		t.progress <- ""
		<-t.waiter
	}
	return nil
}

func (t *Task) Finalize() error {
	if !t.finalized {
		t.progress <- ""
		<-t.waiter
	}
	close(t.kill)
	close(t.progress)
	close(t.waiter)
	close(t.in)
	return nil
}

func (t *Task) Kill() {
	if !t.finalized {
		close(t.kill)
	}
}

func (t *Task) execute() {

	// internal func
	CheckKill := func() bool {
		select {
		case <-t.kill:
			return true
		default:
			return false
		}
	}

	bufout := bytes.NewBuffer(make([]byte, 0))
	buferr := bytes.NewBuffer(make([]byte, 0))
	exitOut := make(chan string, 1)
	exitIn := make(chan string, 1)
	exitErr := make(chan string, 1)
	exitExecute := make(chan string, 1)

	outputKilled := false
	var mutex sync.Mutex

	defer func() {
		bufout.Reset()
		buferr.Reset()
		close(exitIn)
	}()

	stdin, _ := t.cmd.StdinPipe()
	t.cmd.Stdout = bufout
	t.cmd.Stderr = buferr

	WaitForOutputKill := func() {
		mutex.Lock()
		defer mutex.Unlock()
		if !outputKilled {
			outputKilled = true
			exitOut <- "finalize output listener"
			exitErr <- "finalize error listener"
			exitExecute <- "finalize execute end listener"
			time.Sleep(50 * time.Millisecond)
			for _, exist := <-exitOut; exist; _, exist = <-exitOut {
			}
			for _, exist := <-exitErr; exist; _, exist = <-exitErr {
			}
			for _, exist := <-exitExecute; exist; _, exist = <-exitExecute {
			}
		}
	}

	// check kill
	if exit := CheckKill(); exit {
		WaitForOutputKill()
		t.finalized = true
		return
	}

	// call initialize ended
	t.waiter <- Successfully

	// wait for next
	for isContinue := true; isContinue; {
		select {
		case <-t.kill:
			WaitForOutputKill()
			t.finalized = true
			return
		case <-t.progress:
			isContinue = false
			break
		}
	}

	if err := t.cmd.Start(); err != nil {
		// todo: error manage
		return
	}

	// input listener
	go func() {
		var finalizeStr string
		defer func() { exitIn <- finalizeStr }()
		for {
			select {
			case data := <-exitIn:
				finalizeStr = data
			case input := <-t.in:
				if bytes, err := t.fromStr(input.Data()); err != nil {
				} else {
					stdin.Write(bytes)
				}
			}
		}
	}()

	// output listener
	go func() {
		defer close(exitOut)
		preLen := int(0)
		isExit := false

		for {
			select {
			case <-exitOut:
				isExit = true
			default:
				crtLen := bufout.Len()

				if crtLen > preLen && bufout.Bytes()[crtLen-1] != 0 {
					if str, err := t.toStr(bufout.Bytes()[preLen:crtLen]); err != nil {
					} else {
						t.sendOut(str)
						preLen = crtLen
						time.Sleep(time.Duration(t.SyncMs * int(time.Millisecond)))
					}
				} else if isExit {
					return
				}
			}
		}
	}()

	// error listener
	go func() {
		defer close(exitErr)
		preLen := int(0)
		isExit := false

		for {
			select {
			case <-exitErr:
				isExit = true
			default:
				crtLen := buferr.Len()

				if crtLen > preLen && bufout.Bytes()[crtLen-1] != 0 {
					if str, err := t.toStr(bufout.Bytes()[preLen:crtLen]); err != nil {
					} else {
						t.sendErr(str)
						preLen = crtLen
						time.Sleep(time.Duration(t.SyncMs * int(time.Millisecond)))
					}
				} else if isExit {
					return
				}
			}
		}
	}()

	// proccess killer listener
	go func() {
		if err := t.cmd.Wait(); err != nil {
			// todo: error manage
		}
		close(exitExecute)
	}()

	for isContinue := true; isContinue; {
		select {
		case <-t.kill:
			t.cancel()
			t.finalized = true
			break
		case <-exitExecute:
			isContinue = false
			break
		}
	}

	WaitForOutputKill()

	t.waiter <- Successfully
	if !t.finalized {
		<-t.progress
	}

	defer func() { t.waiter <- Successfully }()
	return
}

func (t *Task) sendOut(msg string) {
	for _, ch := range t.out {
		ch <- &task.Packet{
			SendFromName: t.Name,
			DataStr:      msg,
			Color:        t.Color,
		}
	}
}

func (t *Task) sendErr(msg string) {
	for _, ch := range t.err {
		ch <- &task.Packet{
			SendFromName: t.Name,
			DataStr:      msg,
			Color:        t.Color,
		}
	}
}
