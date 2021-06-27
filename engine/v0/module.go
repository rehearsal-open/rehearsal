package v0

import (
	"github.com/rehearsal-open/rehearsal/entity"
	"github.com/rehearsal-open/rehearsal/task"
)

type RehearsalEngine struct {
	*entity.Conf
	tasks map[string]task.Task
}
