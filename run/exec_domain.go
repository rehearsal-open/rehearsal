package run

import (
	"os/exec"
)

type IOExpression struct {
	fromId int
	data   string
	err    IOErrorProps
}

type IOErrorProps struct {
	err      error
	priority int
}

type IOErrorPriority int

const (
	ErrorMustKill IOErrorPriority = iota
	ErrorWarning
	ErrorCaution
	ErrorInfomation
)

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
