// task/impl/net/task.go
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

package net

import (
	"net"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/task/based"
)

func (t *__task) Init(args based.MainFuncArguments) error {

	if conn, err := net.DialTimeout(t.Entity().Kind, t.Address, t.timeout); err != nil {
		return errors.WithStack(err)
	} else {
		t.Conn = conn
	}
	return nil
}