package exec

import (
	"bytes"
	"os/exec"
	"testing"
)

const (
	python1 string = "./../../test/python/python1.py"
	python2 string = "./../../test/python/python2.py"
)

func TestExecFlow(t *testing.T) {

	cmd := exec.Command("python", python1)
	data := make(chan string)
	// out, err := cmd.StdoutPipe()

	// if err != nil {
	// 	t.Log("cmd.StdoutPipe() errored: ", err)
	// 	t.Fail()
	// }

	var buf bytes.Buffer
	cmd.Stdout = &buf

	if err := cmd.Start(); err != nil {
		t.Log("cmd.Start() errored: ", err)
		t.Fail()
	}

	go func() {
		t.Log("output getter initialized.")
		beforeLength := int(0)
		isExit := false
		defer func() { data <- "killed process" }()

		for {
			select {
			case data := <-data:
				t.Log(data)
				isExit = true
			default:

				len := buf.Len()

				if len > beforeLength {
					bytes := buf.Bytes()[beforeLength:len]
					str := string(bytes)
					t.Log(str)
					beforeLength = len
				} else if isExit {
					t.Log(beforeLength, ", ", len)
					return
				}
			}
		}
	}()

	cmd.Wait()
	func() { data <- "process ended" }()
	returnMsg := <-data
	t.Log(returnMsg)
	t.Log("python 1 process ended.")
}
