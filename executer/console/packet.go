package console

import (
	"github.com/rehearsal-open/rehearsal/executer"
)

func (p *Packet) SendFrom() string {
	return p.sendFrom
}

func (p *Packet) StringData() string {
	return p.data
}

func (p *Packet) ErrorInfo() executer.ErrorPriority {
	return p.priority
}
