package out

import (
	"github.com/rehearsal-open/rehearsal/engine"
	"github.com/rehearsal-open/rehearsal/logger"
	"github.com/rehearsal-open/rehearsal/packet"
)

type Task struct {
	engine *engine.RehearsalEngine
	logger *logger.Logger
	in     chan packet.Packet
}

func (t *Task) AssignEngine(e *engine.RehearsalEngine) error {
	t.engine = e
	t.in = make(chan packet.Packet)
}
