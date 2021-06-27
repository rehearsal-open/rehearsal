package task

import (
	color "github.com/rehearsal-open/rehearsal/entity/cli-color"
)

type Packet struct {
	SendFromName string
	DataStr      string
	Color        color.CliColor
}

func (p *Packet) SendFrom() string {
	return p.SendFromName
}

func (p *Packet) Data() string {
	return p.DataStr
}

func (p *Packet) ForeColor() color.CliColor {
	return p.Color
}

func (p *Packet) BackColor() color.CliColor {
	return color.Default
}

func (p *Packet) ConsoleOut() string {
	return color.Fore(p.Color) + p.DataStr
}
