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
	"errors"
	"regexp"
)

const (
	// for user defined name (tasks, phases, and so on...)
	userDefinedName = "[a-zA-Z][a-zA-Z0-9_-]*"
	// for system defined name (tasks, phases, and so on...)
	systemDefinedName = "__[a-zA-Z][a-zA-Z0-9_-]*"
	// for user and system defined name (tasks, phases, and so on...)
	definedName = "(__)?[a-zA-Z][a-zA-Z0-9_-]*"
	// for task's fullname, including phase and task's name
	taskFullName = definedName + "::" + definedName
	// for user defined task's fullname, including phase and task user's defined name
	userTaskFullName = userDefinedName + "::" + userDefinedName
)

type __Regexp struct {
	userDefinedName   *regexp.Regexp // for user defined name(tasks, phases, and so on): "[a-zA-Z][a-zA-Z0-9_-]*"
	systemDefinedName *regexp.Regexp // for system defined name(tasks, phases, and so on): "__[a-zA-Z][a-zA-Z0-9_-]*"
	definedName       *regexp.Regexp // for user and system name(tasks, phases, and so on): "(__)?[a-zA-Z][a-zA-Z0-9_-]*"
	taskFullName      *regexp.Regexp // for task's fullname: "(__)?[a-zA-Z][a-zA-Z0-9_-]*::(__)?[a-zA-Z][a-zA-Z0-9_-]*"
	userTaskFullName  *regexp.Regexp // for user defined task's fullname: "[a-zA-Z][a-zA-Z0-9_-]*::[a-zA-Z][a-zA-Z0-9_-]*"
}

var Regexp __Regexp = __Regexp{
	userDefinedName:   regexp.MustCompile(userDefinedName),
	systemDefinedName: regexp.MustCompile(systemDefinedName),
	definedName:       regexp.MustCompile(definedName),
	taskFullName:      regexp.MustCompile(taskFullName),
	userTaskFullName:  regexp.MustCompile(userTaskFullName),
}

// regular expression for user defined name(tasks, phases, and so on): "[a-zA-Z][a-zA-Z0-9]*"
func (r *__Regexp) UserDefinedName() *regexp.Regexp { return r.userDefinedName }

// for system defined name(tasks, phases, and so on): "__[a-zA-Z][a-zA-Z0-9]*"
func (r *__Regexp) SystemDefinedName() *regexp.Regexp { return r.systemDefinedName }

// for user and system name(tasks, phases, and so on): "([a-zA-Z]|__[a-zA-Z])[a-zA-Z0-9]*"
func (r *__Regexp) DefinedName() *regexp.Regexp { return r.definedName }

// for task's fullname: "([a-zA-Z]|__[a-zA-Z])[a-zA-Z0-9]*::([a-zA-Z]|__[a-zA-Z])[a-zA-Z0-9]*"
func (r *__Regexp) TaskFullName() *regexp.Regexp { return r.taskFullName }

// for user defined task's fullname: "[a-zA-Z][a-zA-Z0-9]*::[a-zA-Z][a-zA-Z0-9]*"
func (r *__Regexp) UserTaskFullName() *regexp.Regexp { return r.userTaskFullName }

// string value satisfied entitities.Regexp.DefinedName regular expression
type DefinedName string

// check value satisfied regular expression and assign it
func (defined DefinedName) Assign(name string) error {
	if found := Regexp.definedName.FindAllStringIndex(name, -1); len(found) == 1 {
		if found[0][0] == 0 && found[0][1] == len(name) {
			defined = DefinedName(name)
			return nil
		}
	}
	return errors.New(name + " is invalid (not satisfied DefinedName regular expression)")
}

type TaskFullName string

func MakeTaskFullName(phaseName DefinedName, taskName DefinedName) TaskFullName {
	return TaskFullName(phaseName + "::" + taskName)
}
