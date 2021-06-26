package executer_test

import (
	"testing"

	"github.com/rehearsal-open/rehearsal/executer"
	"github.com/rehearsal-open/rehearsal/executer/console"
)

const (
	pythonInOut string = "./../test/python/python1.py"
	pythonOut   string = "./../test/python/python2.py"
)

func TestPythonToPython(t *testing.T) {
	python1 := console.ConsoleExecuteMaker("python", pythonOut)
	python2 := console.ConsoleExecuteMaker("python", pythonInOut)
	manager := executer.ExecuteManagerMaker()

	manager.Append("out", &executer.ExecuteItem{
		Exec: python1,
	})

	manager.Append("inout", &executer.ExecuteItem{
		Exec: python2,
	})

	manager.ConnectOut("out", "inout")
	manager.ExecuteAll()
	return
}
