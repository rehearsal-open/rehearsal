package engine

import (
	"github.com/rehearsal-open/rehearsal/entity"
)

type RehearsalEngine interface {
	AssignConfig(conf *entity.Config) error
	Config() *entity.Config
	Run() error
	Kill()
	Finalize()
}
