package exec

import (
	"os/exec"
)

type Packet struct {
	pid      int
	data     string
	priority ErrorPriority
}

type ErrorPriority int

const (
	ErrorInfomation ErrorPriority = 1 + iota
	ErrorCaution
	ErrorWarning
	ErrorFatal
)

type Exec struct {
	cmd *exec.Cmd
	in  chan Packet
	out [](chan Packet)
	err [](chan Packet)
}

func ExecMaker(cmd *exec.Cmd) *Exec {
	result := Exec{
		cmd: cmd,
		in:  make(chan Packet),
		out: make([]chan Packet, 0),
		err: make([]chan Packet, 0),
	}
	return &result
}

func (exec *Exec) AppendOutExec(appended *Exec) {
	exec.out = append(exec.out, appended.in)
}

func (exec *Exec) AppendOutChannel(appended chan Packet) {
	exec.out = append(exec.out, appended)
}

func (exec *Exec) AppendErrExec(appended *Exec) {
	exec.err = append(exec.err, appended.in)
}

func (exec *Exec) AppendErrChannel(appended chan Packet) {
	exec.err = append(exec.err, appended)
}
