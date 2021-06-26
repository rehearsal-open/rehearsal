package console

import (
	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/executer"
)

func (e *Execute) IsRunning() bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.isrun
}

func (e *Execute) IsRunningError() error {
	if isRun := e.IsRunning(); isRun {
		return errors.New("This execute process is running.")
	} else {
		return nil
	}
}

func (e *Execute) SetName(name string) error {
	if err := e.IsRunningError(); err != nil {
		return errors.WithStack(err)
	} else {
		e.mutex.Lock()
		defer e.mutex.Unlock()
		e.name = name
		return nil
	}
}

func (e *Execute) GetName() string {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.name
}

func (e *Execute) AppendOutPacketChannel(reciever chan executer.Packet) error {
	if err := e.IsRunningError(); err != nil {
		return errors.WithStack(err)
	} else {
		e.mutex.Lock()
		defer e.mutex.Unlock()
		e.out = append(e.out, reciever)
		return nil
	}
}

func (e *Execute) AppendErrPacketChannel(reciever chan executer.Packet) error {
	if err := e.IsRunningError(); err != nil {
		return errors.WithStack(err)
	} else {
		e.mutex.Lock()
		defer e.mutex.Unlock()
		e.err = append(e.err, reciever)
		return nil
	}
}

func (e *Execute) GetInputChannel() chan executer.Packet {
	return e.in
}

func (e *Execute) SetConvToString(conv func([]byte) (string, error)) error {
	if err := e.IsRunningError(); err != nil {
		return errors.WithStack(err)
	} else {
		e.mutex.Lock()
		defer e.mutex.Unlock()
		e.toStr = conv
		return nil
	}
}

func (e *Execute) SetConvFromString(conv func(string) ([]byte, error)) error {
	if err := e.IsRunningError(); err != nil {
		return errors.WithStack(err)
	} else {
		e.mutex.Lock()
		defer e.mutex.Unlock()
		e.fromStr = conv
		return nil
	}
}

func (e *Execute) SetKillAll(ch chan string) error {
	if err := e.IsRunningError(); err != nil {
		return errors.WithStack(err)
	} else {
		e.mutex.Lock()
		defer e.mutex.Unlock()
		e.killall = ch
		return nil
	}
}
