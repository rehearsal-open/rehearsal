package console

func (e *Execute) BytesToString(src []byte) (string, error) {
	return e.toStr(src)
}

func (e *Execute) BytesFromString(src string) ([]byte, error) {
	return e.fromStr(src)
}
