package packet

import (
	color "github.com/rehearsal-open/rehearsal/entity/cli-color"
)

type Packet interface {
	SendFrom() string
	Data() string
	ConsoleOut() string
	ForeColor() color.CliColor
	BackColor() color.CliColor
}
