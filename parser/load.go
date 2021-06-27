package parser

import (
	"github.com/pkg/errors"

	"github.com/rehearsal-open/rehearsal/entity"
)

func Load() (*entity.Conf, error) {

	conf := entity.Conf{}
	if err := LoadArgs(&conf); err != nil {
		return nil, errors.WithStack(err)
	} else {
		return &conf, nil
	}
}
