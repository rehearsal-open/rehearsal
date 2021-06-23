package run

import (
	"os/exec"
)

type IOExpression struct {
	fromId int
	data   string
	other  interface{}
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
	Recieve   chan IOExpression
	errSendTo []*(chan IOExpression)
	state     ExecRunOn
}
