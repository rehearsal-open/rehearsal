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
	"os"
	"time"

	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/parser/mapped"
	"github.com/streamwest-1629/convertobject"
)

type (
	Detail struct {
		Path      string   `map-to:"cmd!"`
		Args      []string `map-to:"args"`
		Dir       string   `map-to:"dir"`
		IsWait    bool     `map-to:"wait-stop"`
		Timelimit time.Duration
	}
)

var (
	DirPathDefault string
)

func init() {
	if wd, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		DirPathDefault = wd
	}
}

func (d *Detail) CheckFormat() error {
	return nil
}

func GetDetail(mapping mapped.MappingType, dest *entities.Task) error {
	// TODO: WRITE IT

	detail := &Detail{
		IsWait: true,
		Dir:    DirPathDefault,
		Args:   []string{},
	}

	if err := convertobject.DirectConvert(mapping, detail); err != nil {
		return err
	} else {
		dest.IsWait, dest.Detail = detail.IsWait, detail
	}

	return nil
}

func (d *Detail) String() string {
	return ""
}
