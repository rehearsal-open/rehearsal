package cli

import (
	"testing"
)

func TestCliTextColor(t *testing.T) {
	str := CliStr{isBold: true, BackColor: Red, TextColor: Yellow, Text: "警告"}
	t.Log(str.GetString())
}
