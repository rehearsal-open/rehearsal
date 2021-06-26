package executer

import (
	"github.com/pkg/errors"
)

func (manager *ExecuteManager) GetMaxNameLen() int {
	return manager.maxNameLen
}

func (manager *ExecuteManager) GetExecute(name string) (*ExecuteItem, error) {
	if executer, exist := manager.executes[name]; exist == false {
		return nil, errors.New("Cannot find execute task or not registered yet: " + name)
	} else {
		return executer, nil
	}
}

func (manager *ExecuteManager) Append(name string, item *ExecuteItem) error {
	if _, exist := manager.executes[name]; exist {
		return errors.New("Already registered execute task: " + name)
	} else {
		item.Exec.SetName(name)
		manager.executes[name] = item
		if len(name) > manager.maxNameLen {
			manager.maxNameLen = len(name)
		}
		return nil
	}
}

func (manager *ExecuteManager) ConnectOut(sendName string, recieveName string) error {
	if send, err := manager.GetExecute(sendName); err != nil {
		return errors.WithStack(err)
	} else if recieve, err := manager.GetExecute(recieveName); err != nil {
		return errors.WithStack(err)
	} else if err := send.Exec.AppendOutPacketChannel(recieve.Exec.GetInputChannel()); err != nil {
		return errors.WithStack(err)
	} else {
		return nil
	}
}

func (manager *ExecuteManager) ConnectErr(sendName string, recieveName string) error {
	if send, err := manager.GetExecute(sendName); err != nil {
		return errors.WithStack(err)
	} else if recieve, err := manager.GetExecute(recieveName); err != nil {
		return errors.WithStack(err)
	} else if err := send.Exec.AppendErrPacketChannel(recieve.Exec.GetInputChannel()); err != nil {
		return errors.WithStack(err)
	} else {
		return nil
	}
}

func (manager *ExecuteManager) ConnectOutChannel(sendName string, recieve chan Packet) error {
	if send, err := manager.GetExecute(sendName); err != nil {
		return errors.WithStack(err)
	} else if err := send.Exec.AppendOutPacketChannel(recieve); err != nil {
		return errors.WithStack(err)
	} else {
		return nil
	}
}

func (manager *ExecuteManager) ConnectErrChannel(sendName string, recieve chan Packet) error {
	if send, err := manager.GetExecute(sendName); err != nil {
		return errors.WithStack(err)
	} else if err := send.Exec.AppendErrPacketChannel(recieve); err != nil {
		return errors.WithStack(err)
	} else {
		return nil
	}
}
