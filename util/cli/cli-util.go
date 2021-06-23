package cli

import (
	"fmt"
)

type Color int

const ColorDef Color = -1
const (
	Black Color = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Syan
	White
)

type CliStr struct {
	TextColor Color
	BackColor Color
	Text      string
	isBold    bool
	isItalic  bool
	isULine   bool
}

func (cli *CliStr) GetString() (result string) {

	result = ""
	if cli.isBold {
		result += "\x1b[1m"
	}
	if cli.isItalic {
		result += "\x1b[3m"
	}
	if cli.isULine {
		result += "\x1b[4m"
	}
	if cli.BackColor > 0 {
		result += fmt.Sprintf("\x1b[%dm", 40+cli.BackColor)
	}
	if cli.TextColor > 0 {
		result += fmt.Sprintf("\x1b[%dm", 30+cli.TextColor)
	}
	result += fmt.Sprintf("%s\x1b[0m", cli.Text)
	return
}
