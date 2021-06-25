package run

import (
	"github.com/pkg/errors"

	"github.com/rehearsal-open/rehearsal/run/exec"
	"github.com/rehearsal-open/rehearsal/util/cli"
)

type ExecuteOption struct {
	Exec  *exec.Exec
	Title string
	color cli.Color
}

type ExecutesManager struct {
	executers map[string]ExecuteOption
}

func ExecutesMaker() *ExecutesManager {
	return &ExecutesManager{
		executers: make(map[string]ExecuteOption),
	}
}

func (manager *ExecutesManager) AppendExecute(option *ExecuteOption) error {
	if _, exist := manager.executers[option.Title]; exist {
		return errors.New("Already registered execute task: " + option.Title)
	}
	manager.executers[option.Title] = *option

	return nil
}

func (manager *ExecutesManager) GetExecuter(name string) (*ExecuteOption, error) {
	executer, exist := manager.executers[name]
	if !exist {
		return nil, errors.New("Cannot find execute task or not registered yet: " + name)
	}
	return &executer, nil
}

func (manager *ExecutesManager) AppendConnector(senderExecName string, recieverExecName string) error {

	sender, err := manager.GetExecuter(senderExecName)
	if err != nil {
		return errors.WithStack(err)
	}
	reciever, err := manager.GetExecuter(recieverExecName)
	if err != nil {
		return errors.WithStack(err)
	}

	sender.Exec.AppendOutExec(reciever.Exec)
	return nil
}

func (manager *ExecutesManager) AppendErrConnector(senderExecName string, recieverExecName string) error {

	sender, err := manager.GetExecuter(senderExecName)
	if err != nil {
		return errors.WithStack(err)
	}
	reciever, err := manager.GetExecuter(recieverExecName)
	if err != nil {
		return errors.WithStack(err)
	}
	sender.Exec.AppendErrExec(reciever.Exec)
	return nil
}

func (manager *ExecutesManager) AppendOutChannel(senderExecName string, ch chan exec.Packet) error {
	sender, err := manager.GetExecuter(senderExecName)
	if err != nil {
		return errors.WithStack(err)
	}
	sender.Exec.AppendOutChannel(ch)
	return nil
}

func (manager *ExecutesManager) AppendErrChannel(senderExecName string, ch chan exec.Packet) error {
	sender, err := manager.GetExecuter(senderExecName)
	if err != nil {
		return errors.WithStack(err)
	}
	sender.Exec.AppendErrChannel(ch)
	return nil
}
