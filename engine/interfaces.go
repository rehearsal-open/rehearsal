package engine

import (
	"github.com/rehearsal-open/rehearsal/entity"
)

type RehearsalEngine interface {
	AssignConfig(config *entity.Conf) error
	Execute() error
	Kill() error
}
