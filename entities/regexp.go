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

	"github.com/streamwest-1629/guarantee_str"
)

const (
	// for user defined name (tasks, phases, and so on...)
	userDefinedNameRegex = "[a-zA-Z][a-zA-Z0-9_-]*"
	// for system defined name (tasks, phases, and so on...)
	systemDefinedNameRegex = "__[a-zA-Z][a-zA-Z0-9_-]*"
	// for user and system defined name (tasks, phases, and so on...)
	definedNameRegex = "(__)?[a-zA-Z][a-zA-Z0-9_-]*"
	// for task's fullname, including phase and task's name
	taskFullNameRegex = definedNameRegex + "::" + definedNameRegex
	// for user defined task's fullname, including phase and task user's defined name
	userTaskFullNameRegex = userDefinedNameRegex + "::" + userDefinedNameRegex
)

var (
	// for user defined name(tasks, phases, and so on): "[a-zA-Z][a-zA-Z0-9_-]*"
	userDefinedNameFilter = guarantee_str.MakeRegexpFilter(regexp.MustCompile(userDefinedNameRegex))
	// for system defined name(tasks, phases, and so on): "__[a-zA-Z][a-zA-Z0-9_-]*"
	systemDefinedNameFilter = guarantee_str.MakeRegexpFilter(regexp.MustCompile(systemDefinedNameRegex))
	// for user and system name(tasks, phases, and so on): "(__)?[a-zA-Z][a-zA-Z0-9_-]*"
	definedNameFilter = guarantee_str.MakeRegexpFilter((regexp.MustCompile(definedNameRegex)))
	// for task's fullname: "(__)?[a-zA-Z][a-zA-Z0-9_-]*::(__)?[a-zA-Z][a-zA-Z0-9_-]*"
	taskFullNameFilter = guarantee_str.MakeRegexpFilter(regexp.MustCompile(taskFullNameRegex))
	// for user defined task's fullname: "[a-zA-Z][a-zA-Z0-9_-]*::[a-zA-Z][a-zA-Z0-9_-]*"
	userTaskFullNameFilter = guarantee_str.MakeRegexpFilter(regexp.MustCompile(userTaskFullNameRegex))
)

// to make internal embed structure
type internalGuarantee struct{ *guarantee_str.GuaranteeStr }

type DefinedName struct{ *internalGuarantee }

func MakeDefinedName(name string) (*DefinedName, error) {
	if guarantee, err := definedNameFilter.MakeGuarantee(name); err != nil {
		return nil, err
	} else {
		return &DefinedName{
			internalGuarantee: &internalGuarantee{
				GuaranteeStr: guarantee,
			},
		}, nil
	}
}

func (name *DefinedName) IsSystemName() bool {
	_, err := systemDefinedNameFilter.ChangeGuarantee(name.GuaranteeStr)
	return err == nil
}

// task's fullname, including phase name and task name
type TaskFullName struct{ *internalGuarantee }

func MakeTaskFullName(phaseName, taskName string) (*TaskFullName, error) {
	if guarantee, err := definedNameFilter.MakeGuarantee(phaseName + "::" + taskName); err != nil {
		return nil, err
	} else {
		return &TaskFullName{
			internalGuarantee: &internalGuarantee{
				GuaranteeStr: guarantee,
			},
		}, nil
	}
}
