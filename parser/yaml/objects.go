// parser/yaml/objects.go
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

package yaml

import (
	"os"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/parser"
	"github.com/rehearsal-open/rehearsal/parser/mapped"
	"gopkg.in/yaml.v3"
)

type (
	Parser struct {
		Path string
		parser.DetailMaker
	}
)

func (p *Parser) Parse(init parser.EnvConfig, dest *entities.Rehearsal) error {

	var (
		mapping = parser.MappingType{}
		file    *os.File
	)

	if f, err := os.Open(p.Path); err != nil {
		return errors.WithMessage(err, "cannot open yaml config file")
	} else {
		file = f
	}

	defer file.Close()
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&mapping); err != nil {
		return errors.WithMessage(err, "cannot read yaml config file")
	}

	parser := mapped.Parser{
		Mapped:      mapping,
		DetailMaker: p.DetailMaker,
	}

	err := parser.Parse(init, dest)
	return errors.WithMessage(err, "cannot load config data")
}
