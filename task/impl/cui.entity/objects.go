// task/impl/cui.entity/objects.go
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

package cui

import (
	"time"

	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/parser/mapped"
)

type (
	Detail struct {
		Path      string
		Args      []string
		Dir       string
		Timelimit time.Duration
	}
)

func (d *Detail) CheckFormat() error {
	return nil
}

func GetDetail(mapping mapped.MappingType, dest *entities.Task) error {
	// TODO: WRITE IT
	return nil
}

func (d *Detail) String() string {
	return ""
}
