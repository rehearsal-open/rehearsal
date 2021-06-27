package packet

type Packet interface {
	SendFrom() string
	Data() string
	ConsoleOut() string
}
