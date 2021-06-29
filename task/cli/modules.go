package cli

import (
	"github.com/rehearsal-open/rehearsal/engine"
	"github.com/rehearsal-open/rehearsal/entity"
	"github.com/rehearsal-open/rehearsal/packet"
)

type Task struct {
	engine   *engine.RehearsalEngine
	taskConf *entity.TaskConfig
	in       chan packet.Packet
	out      []chan packet.Packet
	err      []chan packet.Packet
	toStr    func([]byte) (string, error)
	fromStr  func([]string) ([]byte, error)
}
