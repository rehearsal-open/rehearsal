package task

import (
	color "github.com/rehearsal-open/rehearsal/entity/cli-color"
)

type Packet struct {
	sendFrom string
	data     string
	color    color.CliColor
}

func (p *Packet) SendFrom() string {
	return p.sendFrom
}

func (p *Packet) Data() string {
	return p.data
}

func (p *Packet) ForeColor() color.CliColor {
	return p.color
}

func (p *Packet) BackColor() color.CliColor {
	return color.Default
}

func (p *Packet) ConsoleOut() string {
	return color.Fore(p.color) + p.data
}
