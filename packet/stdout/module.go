package stdout

type Packet struct {
	Name         string
	cliDecorated string
	Data         string
}

func (p *Packet) SendFrom() string {
	return p.Name
}

func (p *Packet) GetString() string {
	return p.Data
}

func (p *Packet) CLIView() string {
	return p.Data
}
