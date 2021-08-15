// entities/constants.go
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

package entities

type (
	RunningState int
	TaskElement  int
)

const (
	Main TaskElement = iota
	StdIn
	StdOut
	StdErr
	numTaskElement
	NumTaskElement int = int(numTaskElement)
)

const (
	Waiting RunningState = iota
	Running
	Closed
	Finalized
)

func (t TaskElement) String() string {
	return [...]string{
		"main task",
		"standard input",
		"standard output",
		"standard error output",
	}[t]
}
