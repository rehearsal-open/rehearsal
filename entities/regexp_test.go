// entities/regexp_test.go
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

package entities_test

import (
	"testing"

	"github.com/rehearsal-open/rehearsal/entities"
)

func TestTaskName(t *testing.T) {

	testStr := [...]string{
		"aiueo",
		"a1ue0",
		"01234",
		"__012",
		"__aab",
	}

	for _, test := range testStr {

		if name, err := entities.MakeDefinedName(test); err != nil {
			t.Log(err.Error())
		} else {
			if isSystem := name.IsSystemName(); isSystem {
				t.Log("Success: " + name.MustCloneString() + " (system name)")
			} else {
				t.Log("Success: " + name.MustCloneString())
			}
		}
	}
}
