// entities/objects.go
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
	"time"

	"github.com/streamwest-1629/textfilter"
)

type (
	// Defines configuration of rehearsal executing and eachtask default configuration's value.
	Rehearsal struct {
		// Phase name filter used when appended phases.
		phaseFilter textfilter.Filter
		// I/O synchronization interval of each task's default value and system I/O.
		SyncInterval time.Duration
		phases       []*Phase
		nameList     map[string]int
	}

	// Configuration of phase, including some rehearsal tasks.
	Phase struct {
		name         string
		SyncInterval time.Duration
		// Task name filter used when appended phases.
		taskFilter textfilter.Filter
		tasks      []*Task
		nameList   map[string]int
	}

	// Defines configuration of rehearsal task, its lifespan.
	Task struct {
		name string
		Kind string
		// detail configuration indivisually defined by task's kind
		detail TaskDetail
		// I/O synchronization interval.
		SyncInterval time.Duration

		// text Encoding, default value is "utf8"
		Encoding string
	}

	TaskDetail interface {
		CheckFormat() error
	}
)

var (
	taskkindFilter = textfilter.ListMatches(
		"cui",
		"console",
	)
	encodingFilter = textfilter.ListMatches(
		"utf8",
	)
)
