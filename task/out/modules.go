package out

import (
	"time"

	"github.com/rehearsal-open/rehearsal/engine"
	"github.com/rehearsal-open/rehearsal/logger"
	"github.com/rehearsal-open/rehearsal/packet"
)

type Task struct {
	engine      engine.RehearsalEngine
	logger      *logger.Logger
	in          chan packet.Packet
	exitRoutine chan error
	killed      bool
}

func (t *Task) AssignEngine(e engine.RehearsalEngine, name string) error {
	t.engine = e
	t.in = make(chan packet.Packet)
	t.killed = false

	return nil
}

func (t *Task) AssignLogger(l *logger.Logger) error {
	t.logger = l
	return nil
}

func (t *Task) BytesFromString(src string, sendFrom string) ([]byte, error) {
	return []byte(src), nil
}

func (t *Task) In() chan packet.Packet {
	return t.in
}

func (t *Task) RunInit() error {

	t.exitRoutine = make(chan error, 1)
	go t.routine()
	return nil
}

func (t *Task) RunWait() error {
	return nil
}

func (t *Task) Kill() {
	t.exitRoutine <- nil
	t.killed = true
}

func (t *Task) Finalize() {

	if t.killed {
		t.exitRoutine <- nil
	}

	for {
		time.Sleep(10 * time.Millisecond)
		if _, exist := <-t.exitRoutine; !exist {

			close(t.in)
		}
	}
}

func (t *Task) routine() {
	isContinue := true
	for {
		select {
		case packet, exist := <-t.in:
			if exist {
				t.logger.PacketPrint(packet)
			} else {
				t.logger.SystemPrint("packet is channel is closed")
			}
		case <-t.exitRoutine:
			isContinue = false
		default:
			if !isContinue {
				defer close(t.exitRoutine)
				return
			}
		}
	}
}
