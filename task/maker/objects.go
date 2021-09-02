// task/maker/objects.go
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

package maker

import (
	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/frontend"
	"github.com/rehearsal-open/rehearsal/parser"
	"github.com/rehearsal-open/rehearsal/task"
)

type (
	// task's instance maker
	Maker struct {
		frontend.Frontend
		taskMakers map[string]TaskMaker
	}
	TaskMaker interface {
		MakeDetail(frontend frontend.Frontend, def *entities.Rehearsal, src parser.MappingType, dest *entities.Task) error
		MakeTask(entity *entities.Task) (task.Task, error)
	}
	MakerCollection struct {
		MakeDetailFunc func(frontend frontend.Frontend, entity *entities.Rehearsal, src parser.MappingType, dest *entities.Task) error
		MakeTaskFunc   func(entity *entities.Task) (task.Task, error)
	}
)

var (
	ErrCannotSupportKind = errors.New("cannot found task's kind from supported list")
)

func (m *Maker) RegisterMaker(kind string, maker TaskMaker) {
	if m.taskMakers == nil {
		m.taskMakers = make(map[string]TaskMaker)
	}
	m.taskMakers[kind] = maker
}

func (m *Maker) MakeDetail(def *entities.Rehearsal, src parser.MappingType, dest *entities.Task) error {
	if maker, support := m.taskMakers[dest.Kind]; !support {
		return ErrCannotSupportKind
	} else {
		return errors.WithStack(maker.MakeDetail(m.Frontend, def, src, dest))
	}
}

func (m *Maker) MakeTask(entity *entities.Task) (task.Task, error) {
	if maker, support := m.taskMakers[entity.Kind]; !support {
		return nil, ErrCannotSupportKind
	} else if task, err := maker.MakeTask(entity); err != nil {
		return nil, errors.WithStack(err)
	} else {
		return task, err
	}
}

func (m *Maker) IsSupportedKind(kind string) bool {
	_, support := m.taskMakers[kind]
	return support
}

func (c *MakerCollection) MakeDetail(frontend frontend.Frontend, def *entities.Rehearsal, src parser.MappingType, dest *entities.Task) error {
	return c.MakeDetailFunc(frontend, def, src, dest)
}

func (c *MakerCollection) MakeTask(entity *entities.Task) (task.Task, error) {
	return c.MakeTaskFunc(entity)
}
