package out

import (
	"github.com/rehearsal-open/rehearsal/engine"
	"github.com/rehearsal-open/rehearsal/logger"
)

type Task struct {
	engine *engine.RehearsalEngine
	logger *logger.Logger
	in     chan Packet
}

func (t *Task) AssignEngine(e engine.RehearsalEngine) error {

}
