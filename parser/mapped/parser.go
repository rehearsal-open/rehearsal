// parser/mapped/parser.go
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
	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/parser"
	"github.com/streamwest-1629/convertobject"
)

var (
	ErrConfFileNotSupported = errors.New("configuration file's version is not supported, supported 0.202109.")
)

const errMsgBase = "cannot parse from map to object because of "

// Parse from map[string]interface{} to structures.
func (p *Parser) Parse(init parser.EnvConfig, dest *entities.Rehearsal) error {

	// initialize device configuration
	if err := init.InitConfig(p.Mapped); err != nil {
		return errors.WithMessage(err, errMsgBase+"device configuration")
	}

	// embed object to parse
	r := Rehearsal{Rehearsal: dest}

	// call to convert type
	if err := convertobject.DirectConvert(p.Mapped, &r); err != nil {
		return errors.WithMessage(err, errMsgBase+"invalid format")
	}

	// initialize shared configuration
	// check version
	switch r.Version {
	case 1.202109: // current supported
		return errors.WithStack(p.Parse202109(&r))
	default:
		return ErrConfFileNotSupported
	}

}

type phaseByIndex []*Phase

func (a phaseByIndex) Len() int           { return len(a) }
func (a phaseByIndex) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a phaseByIndex) Less(i, j int) bool { return a[i].Index < a[j].Index }
