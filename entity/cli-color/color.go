package entity

import (
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type CliColor int

const (
	Default CliColor = iota - 1
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Syan
	White
)

func Fore(col CliColor) string {
	if col < Black || White < col {
		return "\x1b[0m"
	} else {
		return "\x1b[" + strconv.Itoa(30+int(col)) + "m"
	}
}

func Back(col CliColor) string {
	if col < Black || White < col {
		return "\x1b[0m"
	} else {
		return "\x1b[" + strconv.Itoa(40+int(col)) + "m"
	}
}

func FromString(str string) (col CliColor, err error) {

	// Default value
	col = Default
	err = nil

	low := strings.ToLower(str)
	switch low {
	case "default":
		col = Default
	case "black":
		col = Black
	case "red":
		col = Red
	case "green":
		col = Green
	case "yellow":
		col = Yellow
	case "Magenta":
		col = Magenta
	case "Syan":
		col = Syan
	case "White":
		col = White
	default:
		err = errors.New("Unknown color: " + str)
	}

	return
}
