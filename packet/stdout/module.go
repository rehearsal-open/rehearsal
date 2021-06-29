package stdout

type Packet struct {
	sendFrom     string
	cliDecorated string
	data         string
}

func (p *Packet) SendFrom() string {
	return p.sendFrom
}

func (p *Packet) GetString() string {
	return p.data
}

func (p *Packet) CLIView() string {
	return p.cliDecorated
}
