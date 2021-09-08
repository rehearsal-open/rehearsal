// task/impl/serial/objects.go
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

package serial

import (
	"github.com/rehearsal-open/rehearsal/task/based"
	"go.bug.st/serial"
)

type (
	__task struct {
		*based.Task
		serial.Port
		*Detail
		close chan error
	}

	Detail struct {
		PortName string `map-to:"port"`
		BaudRate int    `map-to:"baudrate!"`
		DataBits int    `map-to:"databits"`
	}
)
