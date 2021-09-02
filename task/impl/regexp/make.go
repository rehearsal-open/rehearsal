// task/impl/regexp/make.go
// Copyright (C) 2021 Kasai Koji

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package regexp

import (
	"regexp"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/frontend"
	"github.com/rehearsal-open/rehearsal/parser"
	"github.com/rehearsal-open/rehearsal/task"
	"github.com/rehearsal-open/rehearsal/task/based"
	"github.com/rehearsal-open/rehearsal/task/maker"
	"github.com/streamwest-1629/convertobject"
)

var MakeCollection = maker.MakerCollection{
	MakeDetailFunc: GetDetail,
	MakeTaskFunc:   Make,
}

func GetDetail(_ frontend.Frontend, def *entities.Rehearsal, mapping parser.MappingType, dest *entities.Task) error {

	detail := &Detail{}

	if err := convertobject.DirectConvert(mapping, detail); err != nil {
		return errors.WithStack(err)
	}

	dest.Detail = detail

	return nil
}

func Make(entity *entities.Task) (t task.Task, err error) {

	result := __task{}

	if detail, ok := entity.Detail.(*Detail); !ok {
		panic("invalid detail objects type")
	} else {
		result.detail = detail
	}

	result.Task = based.MakeBasis(entity, &result)

	if reg, err := regexp.Compile(result.detail.MatchesRegexpr); err != nil {
		return nil, errors.WithStack(err)
	} else {
		result.Regexp = reg
	}

	return &result, nil

}
