package v0

import (
	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entity"
	"github.com/rehearsal-open/rehearsal/logger"
)

func (e RehearsalEngine) AssignConfig(conf *entity.Config) error {

	e.logger = logger.Logger{}
	if err := e.logger.AssignConfig(conf); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
