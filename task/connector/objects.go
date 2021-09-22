// task/connector/objects.go
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

package connector

import "github.com/rehearsal-open/rehearsal/entities"

type (

	// interface uses as input writer
	Writer interface {
		Write(elem *entities.Element, bytes []byte, onFinal func())
	}

	// interface uses as output writer
	Outputter interface {
		AppendInput(writer Writer) error
	}

	// interface uses as output writer
	Filter interface {
		Writer
		Outputter
	}
)
