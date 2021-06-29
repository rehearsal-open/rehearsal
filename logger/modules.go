package logger

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/rehearsal-open/rehearsal/entity"
	"github.com/rehearsal-open/rehearsal/packet"
	"github.com/rehearsal-open/rehearsal/packet/systemlog"
)

type Logger struct {
	config        *entity.Config
	maxNameLength int
	exitRoutine   chan error
	packetChannel chan packet.Packet
}

const SystemCall string = "system"

func (l *Logger) AssignConfig(conf *entity.Config) error {

	// assign
	l.config = conf

	// max name length
	l.maxNameLength = len(SystemCall)
	for _, taskConf := range conf.TaskConf {
		if len(taskConf.Name) > l.maxNameLength {
			l.maxNameLength = len(taskConf.Name)
		}
	}

	// make chan
	l.exitRoutine = make(chan error, 1)
	l.packetChannel = make(chan packet.Packet)

	// go routine
	go l.routine()

	// todo: write definition
	return nil
}

func (l *Logger) SystemPrint(msg string) {
	l.packetChannel <- &systemlog.Packet{
		Msg: msg,
	}
}

func (l *Logger) PacketPrint(packet packet.Packet) {
	l.packetChannel <- packet
}

func (l *Logger) routine() {
	isContinue := true
	for {
		select {
		case packet, exist := <-l.packetChannel:
			if exist {
				str := packet.CLIView()
				from := packet.SendFrom()
				lstr := len(str)

				if lstr > 1 && str[lstr-2:] == "\n\r" {
					str = str[0 : lstr-2]
				} else if lstr > 0 && (str[lstr-1] == '\n' || str[lstr-1] == '\r') {
					str = str[0 : lstr-1]
				}

				outputs := from + strings.Repeat(" ", 1+l.maxNameLength-len(from)) + "\t: "
				str = strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(str, "\n\r", "\n"), "\r", "\n"), "\n", "\n"+strings.Repeat(" ", 17+l.maxNameLength)+"\t\t: ")
				// if strings.Contains(str, "\n") || strings.Contains(str, "\r") {
				// 	outputs = fmt.Sprintln(outputs + "(multi lines...)\n" + str)
				// } else {
				outputs = fmt.Sprintln(outputs + str)
				// }

				log.Print(outputs)
			} else {
				log.Print("packet channel is closed")
			}

		case <-l.exitRoutine:
			isContinue = false
		default:
			if !isContinue {
				defer close(l.exitRoutine)
				return
			} else {
				time.Sleep(time.Duration(l.config.SyncMs))
			}
		}
	}
}

func (l *Logger) Finalize() {

	// close chan
	l.exitRoutine <- nil
	for {
		time.Sleep(10 * time.Millisecond)
		if _, exist := <-l.exitRoutine; !exist {
			close(l.packetChannel)
			return
		}
	}
}
