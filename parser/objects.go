// parser/objects.go
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

package parser

import (
	"github.com/rehearsal-open/rehearsal/entities"
)

type (

	// Using mapping type
	MappingType = map[string]interface{}

	// Defines parse object.
	Parser interface {
		Parse(EnvConfig, *entities.Rehearsal) error
	}

	// The interface to initialize environment configurations.
	// It is useful for supporting task's default value.
	EnvConfig interface {
		InitConfig(src MappingType, entity *entities.Rehearsal) error
	}

	// The interface to set task's detail configurations.
	DetailMaker interface {
		MakeDetail(def *entities.Rehearsal, src MappingType, dest *entities.Task) error
	}
)
