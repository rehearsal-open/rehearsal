package executer

// This interface has data about connection bitween process and process.
// This function should be just return variable.
type Packet interface {
	SendFrom() string
	StringData() string
	ErrorInfo() ErrorPriority
}

// This interface has process executation.
type Execute interface {
	SetName(string) error
	// This function should be just return variable.
	GetName() string
	SetKillAll(chan string) error
	AppendOutPacketChannel(chan Packet) error
	AppendErrPacketChannel(chan Packet) error
	// This function should be just return variable.
	GetInputChannel() chan Packet
	ExecuteInitialize() error
	ExecuteWait() error
	ExecuteFinalize() error
	BytesToString([]byte) (string, error)
	BytesFromString(string) ([]byte, error)
	// This function must not return error.
	// Just show whether process is running or not.
	// And, this function using, you should use mutex.
	IsRunning() bool
}
