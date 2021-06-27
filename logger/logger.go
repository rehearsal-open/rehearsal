package logger

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rehearsal-open/rehearsal/entity"
	color "github.com/rehearsal-open/rehearsal/entity/cli-color"
	. "github.com/rehearsal-open/rehearsal/packet"
)

// Logger using packet
type PacketLogger struct {
	*entity.Conf
	queue chan Packet
	stop  chan interface{}
}

// Begin output log listener
// Listener starts as subprocess
func (logger *PacketLogger) ListenStart() {
	go func() {
		defer close(logger.queue)
		defer close(logger.stop)

		for {
			select {
			case que := <-logger.queue:
				logger.PacketLog(que)
				break
			case <-logger.stop:
				return
			}
		}
	}()
}

// Stop output log listener
func (logger *PacketLogger) ListenEnd() {
	logger.stop <- 0
}

// Output log
func (logger *PacketLogger) PacketLog(packet Packet) {

	// Init Variables
	msg, from := packet.Data(), packet.SendFrom()
	spc := 1 + logger.MaxNameLen - len(from)

	// Write log's header
	output := fmt.Sprint(
		color.Back(packet.ForeColor()),
		color.Fore(packet.BackColor()),
		"[", time.Now(), "] | ", from, strings.Repeat(" ", spc), "|", color.Clear)

	// Cut last carriage return
	if msg[len(from)-2] == '\n' && msg[len(from)-1] == '\r' {
		msg = msg[:len(from)-2]
	} else if msg[len(from)-1] == '\n' || msg[len(from)-1] == '\r' {
		msg = msg[:len(from)-1]
	}

	// Write log's body
	if strings.Contains(msg, "\n") || strings.Contains(msg, "\r") {
		output += "[v multi lines v]" + packet.ConsoleOut()
	} else {
		output += packet.ConsoleOut()
	}

	// View console
	log.Println(output)
}
