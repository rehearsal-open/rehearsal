// entities/regexp.go
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

import (
	"regexp"

	"github.com/streamwest-1629/guarantee/filter"
)

const (
	// for user defined name (tasks, phases, and so on...)
	userDefinedNameExpr = "[a-zA-Z][a-zA-Z0-9_-]*"
	// for system defined name (tasks, phases, and so on...)
	systemDefinedNameExpr = "__[a-zA-Z][a-zA-Z0-9_-]*"
	// for user and system defined name (tasks, phases, and so on...)
	definedNameExpr = "(__)?[a-zA-Z][a-zA-Z0-9_-]*"
	// for task's fullname, including phase and task's name
	taskFullNameExpr = definedNameExpr + "::" + definedNameExpr
	// for user defined task's fullname, including phase and task user's defined name
	userTaskFullNameExpr = userDefinedNameExpr + "::" + userDefinedNameExpr
)

var (
	// Filter for user defined name (tasks, phases, and so on...).
	UserDefinedName = filter.RegexpFilter(regexp.MustCompile(userDefinedNameExpr))
	// Filter for system defined name (tasks, phases, and so on...).
	SystemDefiedName = filter.RegexpFilter(regexp.MustCompile(systemDefinedNameExpr))
	// Filter for user and system defined name (tasks, phases, and so on...).
	DefinedName = filter.RegexpFilter(regexp.MustCompile(definedNameExpr))
	// Filter for task's fullname, including phase and task's name.
	TaskFullName = filter.RegexpFilter(regexp.MustCompile(taskFullNameExpr))
	// Filter for user defined task's fullname, including phase and task user's defined name.
	UserTaskFullName = filter.RegexpFilter(regexp.MustCompile(userTaskFullNameExpr))
)
