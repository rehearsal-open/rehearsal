// engine/v0/modules.go
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

package v0

import (
	"github.com/rehearsal-open/rehearsal/entity"
	"github.com/rehearsal-open/rehearsal/logger"
	"github.com/rehearsal-open/rehearsal/task"
)

type RehearsalEngine struct {
	config *entity.Config
	tasks  map[string]task.Task
	logger *logger.Logger
}

func (e *RehearsalEngine) Config() *entity.Config {
	return e.config
}
