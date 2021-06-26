package executer

import "github.com/rehearsal-open/rehearsal/util/cli"

type ErrorPriority int

const (
	ErrorInfomation ErrorPriority = 1 + iota
	ErrorCaution
	ErrorWarning
	ErrorFatal
)

type ExecuteItem struct {
	Exec  Execute
	Color cli.Color
}

type ExecuteManager struct {
	executes   map[string]*ExecuteItem
	maxNameLen int
}

func ExecuteManagerMaker() *ExecuteManager {
	return &ExecuteManager{
		executes:   make(map[string]*ExecuteItem),
		maxNameLen: 0,
	}
}
