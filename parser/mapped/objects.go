// parser/mapped/objects.go
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
	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/parser"
)

type (
	Parser struct {
		parser.DetailMaker
		Mapped parser.MappingType
	}
	Rehearsal struct {
		*entities.Rehearsal `map-to:"<-"`
		Phases              []*Phase `map-to:"phase!"`
	}
	Phase struct {
		Index int    `map-to:"at"`
		Name  string `map-to:"name!"`
		Tasks []Task `map-to:"task!"`
	}
	Task struct {
		*entities.Task `map-to:"<-"`
		UntilPhase     *string            `map-to:"until"` // default is launching phase
		IsWait         *bool              `map-to:"wait"`  // default is true
		Clone          parser.MappingType `map-to:"<-"`
		StdIn          *Element           `map-to:"stdin"`
		StdOut         *Element           `map-to:"stdout"`
		StdErr         *Element           `map-to:"stderr"`
	}

	Element struct {
		*entities.Element `map-to:"<-"`
		SendTo            []string `map-to:"sendto"`
	}
)
