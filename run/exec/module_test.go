package exec

import (
	"os/exec"
	"testing"
)

func TestExecModule(t *testing.T) {

	cmd1 := exec.Command("python", python1)
	cmd2 := exec.Command("python", python2)
	exe1 := ExecMaker(cmd1)
	exe2 := ExecMaker(cmd2)
	logout := make(chan Packet)

	go func(out chan Packet) {
		for {
			select {
			case data := <-out:
				t.Log(data.Data)
			}
		}
	}(logout)

	exe2.AppendOutExec(exe1)
	exe1.AppendOutChannel(logout)

	// go func(in chan Packet) {
	// 	for i := 0; i < 10000; i++ {
	// 		in <- Packet{
	// 			data: (strconv.Itoa(i) + "\n"),
	// 		}
	// 		time.Sleep(time.Millisecond * 5)
	// 	}
	// }(exe1.in)

	go exe1.Execute()
	go exe2.Execute()
	for {
	}
}
