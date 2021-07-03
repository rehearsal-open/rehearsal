// packet/systemlog/module.go
// Copyright (C) 2021  Kasai Koji

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

package systemlog

type Packet struct {
	Msg string
}

func (p *Packet) SendFrom() string {
	return "system"
}

func (p *Packet) GetString() string {
	return p.Msg
}

func (p *Packet) CLIView() string {
	return p.Msg
}
