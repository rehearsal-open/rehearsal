package engine

import (
	"github.com/rehearsal-open/rehearsal/entity"
	"github.com/rehearsal-open/rehearsal/logger"
)

type RehearsalEngine interface {
	AssignConfig(conf *entity.Config) error
	Config() *entity.Config
	Logger() *logger.Logger
	Run() error
	Kill()
	Finalize()
}
