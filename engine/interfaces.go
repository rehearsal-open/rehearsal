package engine

import (
	"github.com/rehearsal-open/rehearsal/entity"
)

type RehearsalEngine interface {
	AssignConfig(conf *entity.Config) error
	Run() error
	Kill() error
}
