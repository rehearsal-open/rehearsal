package run

import (
	"os/exec"
)

type IOExpression struct {
	fromId int
	data   []byte
}

type ExecRunOn int

const (
	ExecRunOnWaiting ExecRunOn = 1 << iota
	ExecRunOnRunning
	ExecRunOnKilled
)

type Exec struct {
	id        int
	cmd       exec.Cmd
	timeoutMs int64
	sendTo    []*(chan IOExpression)
	Receive   chan IOExpression
	errSendTo []*(chan IOExpression)
	state     ExecRunOn
}
