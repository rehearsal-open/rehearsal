package packet

type Packet interface {
	GetString() string
	CLIView() string
	SendFrom() string
}
