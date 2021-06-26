package console

import (
	"bytes"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/executer"
)

const (
	Successfully  string = "Successfully"
	RunningStart  string = "RunningStart"
	FinalizeStart string = "FinalizeStart"
)

func (e *Execute) ExecuteInitialize() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.isrun {
		return errors.New("Already process is running.")
	}

	e.isrun = true
	e.progress = make(chan string)
	e.waiter = make(chan string)
	go e.execute()

	if initialized := <-e.waiter; initialized != Successfully {
		e.isrun = false
		return errors.New(initialized)
	} else {
		return nil
	}
}

func (e *Execute) ExecuteWait() error {
	e.progress <- RunningStart
	if ended := <-e.waiter; ended != Successfully {
		e.isrun = false
		return errors.New(ended)
	} else {
		return nil
	}
}

func (e *Execute) ExecuteFinalize() error {
	e.progress <- FinalizeStart
	defer func() {
		close(e.progress)
		close(e.waiter)
		e.isrun = false
	}()
	if ended := <-e.waiter; ended != Successfully {
		return errors.New(ended)
	} else {
		return nil
	}
}

func (e *Execute) execute() {

	// Initialize start
	bufout := bytes.NewBuffer(make([]byte, 0))
	buferr := bytes.NewBuffer(make([]byte, 0))
	exitOut := make(chan string, 1)
	exitIn := make(chan string, 1)
	exitErr := make(chan string, 1)
	exitKill := make(chan string, 1)
	outputKilled := false
	var outputKilledLock sync.Mutex

	stdin, _ := e.cmd.StdinPipe()
	e.cmd.Stdout = bufout
	e.cmd.Stderr = buferr

	defer func() {
		bufout.Reset()
		buferr.Reset()
		close(exitIn)
	}()

	WaitForOutputKill := func() {
		outputKilledLock.Lock()
		defer outputKilledLock.Unlock()
		if !outputKilled {
			outputKilled = true
			exitOut <- "Finalize output listener."
			exitErr <- "Finalize error listener."
			time.Sleep(50 * time.Millisecond)
			for _, exist := <-exitOut; exist; _, exist = <-exitOut {
			}
			for _, exist := <-exitErr; exist; _, exist = <-exitErr {
			}
		}
	}

	// Initialize end
	e.waiter <- Successfully
	<-e.progress

	// Running and waiting start
	if err := e.cmd.Start(); err != nil {
		e.waiter <- err.Error()
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
			case input := <-e.in:
				if bytes, err := e.fromStr(input.StringData()); err != nil {
					e.killall <- err.Error()
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
					if str, err := e.toStr(bufout.Bytes()[preLen:crtLen]); err != nil {
						e.killall <- err.Error()
					} else {
						e.sendOut(str)
						preLen = crtLen
					}
				} else if isExit {

					return
				}
			}
		}
	}()

	// error output listener
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
					if str, err := e.toStr(bufout.Bytes()[preLen:crtLen]); err != nil {
						e.killall <- err.Error()
					} else {
						e.sendErr(str)
						preLen = crtLen
					}
				} else if isExit {
					return
				}
			}
		}
	}()

	// process killer listener
	go func() {
		for {
			select {
			case kill := <-e.killall:
				e.cancel()
				WaitForOutputKill()
				e.killall <- kill
			case <-exitKill:
				return
			}
		}
	}()

	if err := e.cmd.Wait(); err != nil {
		e.waiter <- err.Error()
		return
	}

	WaitForOutputKill()

	// Running and waiting end
	e.waiter <- Successfully
	<-e.progress

	// Finalizing start

	// Reset input channel
	for isExist := true; isExist; {
		select {
		case <-e.in:
		default:
			isExist = false
		}
	}

	defer func() { e.waiter <- Successfully }()
	return
}

func (e *Execute) sendOut(msg string) {
	for _, ch := range e.out {
		ch <- &Packet{
			sendFrom: e.name,
			data:     msg,
			priority: executer.ErrorInfomation,
		}
	}
}

func (e *Execute) sendErr(msg string) {
	for _, ch := range e.err {
		ch <- &Packet{
			sendFrom: e.name,
			data:     msg,
			priority: executer.ErrorWarning,
		}
	}
}
