package exec

import (
	"bytes"
	"io"
)

func (exec *Exec) Execute() {

	bufout := bytes.NewBuffer(make([]byte, 0))
	buferr := bytes.NewBuffer(make([]byte, 0))
	exitOut := make(chan string, 1)
	exitIn := make(chan string, 1)
	// exitErr := make(chan string, 1)

	exec.cmd.Stdout = bufout
	exec.cmd.Stderr = buferr
	stdin, _ := exec.cmd.StdinPipe()

	if err := exec.cmd.Start(); err != nil {
		exec.sendErr(err.Error(), ErrorFatal)
	}

	// input listener
	go func() {
		var finalizeStr string
		defer func() { exitIn <- finalizeStr }()
		for {
			select {
			case data := <-exitIn:
				finalizeStr = data
			case input := <-exec.in:
				io.WriteString(stdin, input.data)
			}
		}
	}()

	// output listener
	go func() {
		var finalizeStr string
		defer func() { exitOut <- finalizeStr }()
		preLen := int(0)
		isExit := false

		for {
			select {
			case data := <-exitOut:
				finalizeStr = data
				isExit = true
			default:
				crtLen := bufout.Len()

				if crtLen > preLen && bufout.Bytes()[crtLen-1] != 0 {
					bytes := bufout.Bytes()[preLen:crtLen]
					exec.sendOut(string(bytes))
					preLen = crtLen
					break
				} else if isExit {
					return
				}
			}
		}
	}()

	if err := exec.cmd.Wait(); err != nil {
		exec.sendErr("ProcessErrored: "+err.Error(), ErrorWarning)
	}
	exitOut <- "Finalize output listener."
	exitIn <- "Finalize input listener."
}

func (exec *Exec) sendErr(msg string, priority ErrorPriority) {
	for _, ch := range exec.err {
		ch <- Packet{
			pid:      exec.cmd.Process.Pid,
			data:     msg,
			priority: priority,
		}
	}
}

func (exec *Exec) sendOut(msg string) {
	for _, ch := range exec.out {
		ch <- Packet{
			pid:      exec.cmd.Process.Pid,
			data:     msg,
			priority: ErrorInfomation,
		}
	}
}
