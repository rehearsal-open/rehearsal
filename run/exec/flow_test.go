package exec

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"
)

const (
	python1 string = "./../../test/python/python1.py"
	python2 string = "./../../test/python/python2.py"
)

func TestExecFlow(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, "python", python1)
	data := make(chan string)
	// out, err := cmd.StdoutPipe()

	// if err != nil {
	// 	t.Log("cmd.StdoutPipe() errored: ", err)
	// 	t.Fail()
	// }

	bufout := bytes.NewBuffer(make([]byte, 0))
	// bufout.Grow(100000)
	// bufin := bytes.NewBuffer(make([]byte, 0))
	cmd.Stdout = bufout
	// cmd.Stdin = bufin
	stdin, _ := cmd.StdinPipe()

	input := make(chan string)
	output := make(chan string)

	if err := cmd.Start(); err != nil {
		t.Log("cmd.Start() errored: ", err)
		t.Fail()
	}

	go func(enter chan string) {
		for i := 0; i < 10000; i++ {
			enter <- (strconv.Itoa(i) + "\n")
			time.Sleep(time.Millisecond + 5)
		}
	}(input)

	go func() {
		t.Log("input setter initialized.")
		for {
			select {
			case input := <-input:
				io.WriteString(stdin, input)

			case output := <-output:
				t.Log(output)
			}

			// bufin.WriteString(strconv.Itoa(i) + "\n")
		}
	}()

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

				len := bufout.Len()

				if len > beforeLength && bufout.Bytes()[len-1] != 0 {
					bytes := bufout.Bytes()[beforeLength:len]
					// t.Log("bytes is valid")
					str := string(bytes)
					output <- str
					// t.Log(str, "(", beforeLength, ", ", len, ")")
					beforeLength = len

					time.Sleep(time.Millisecond)

					break
					// for {

					// 	if r, _ := utf8.DecodeLastRune(bytes); r != utf8.RuneError {

					// 	} else if len == beforeLength {
					// 		break
					// 	} else {
					// 		len--
					// 		t.Log("bytes is invalid")
					// 	}
					// }

				} else if isExit {
					os.WriteFile("test_out.txt", bufout.Bytes(), 0666)
					t.Log(beforeLength, ", ", len)
					return
				}
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		t.Log("process errored: " + err.Error())
	}
	func() { data <- "process ended" }()
	returnMsg := <-data
	t.Log(returnMsg)
	t.Log("python 1 process ended.")
}
