// rehearsal-cli/help.go
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

package main

import "fmt"

const (
	helpDefault = `
    -- Usage 'rehearsal' --
rehearsal-cli [command] (args...)

[command] REQUIRED
  - run     : Execute tasks with config file. See also 'rehearsal-cli run help'
  - version : Show rehearsal version.
`
	helpRun = `
    -- Usage 'rehearsal-cli run' --
rehearsal-cli run [config-filepath]

[config-filepath] REQUIRED
   You need to required config file's path.
  In rehearsal, '*.yaml' is supported config file's format type.
	`
)

func PutHelpDefault() {
	fmt.Println(helpDefault)
}
