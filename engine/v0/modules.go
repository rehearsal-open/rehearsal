package v0

import (
	"github.com/rehearsal-open/rehearsal/entity"
	"github.com/rehearsal-open/rehearsal/logger"
	"github.com/rehearsal-open/rehearsal/task"
)

type RehearsalEngine struct {
	config *entity.Config
	tasks  map[string]task.Task
	logger *logger.Logger
}

func (e *RehearsalEngine) Config() *entity.Config {
	return e.config
}

func (e *RehearsalEngine) Logger() *logger.Logger {
	return e.logger
}
