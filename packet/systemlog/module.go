package systemlog

type Packet struct {
	Msg string
}

func (p *Packet) SendFrom() string {
	return "system"
}

func (p *Packet) GetString() string {
	return p.Msg
}

func (p *Packet) CLIView() string {
	return p.Msg
}
