// parser/mapped/objects_test.go
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

package mapped

import (
	"fmt"
	"testing"

	"github.com/streamwest-1629/convertobject"
)

func TestParse(t *testing.T) {
	val := MappingType{
		"version": 0.202109,
		"phase": []interface{}{
			map[string]interface{}{
				"name": "begin",
				"task": []interface{}{
					map[string]interface{}{
						"name": "task-1",
					},
				},
			},
		},
	}

	r := Rehearsal{}

	if err := convertobject.DirectConvert(map[string]interface{}(val), &r); err != nil {
		panic(err.Error())
	}

	fmt.Print(r)
}
