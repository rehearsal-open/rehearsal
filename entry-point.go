// entry-point.go
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

package main

import (
	"fmt"
	"log"

	v0 "github.com/rehearsal-open/rehearsal/engine/v0"
	"github.com/rehearsal-open/rehearsal/entity"
	"github.com/rehearsal-open/rehearsal/load"
)

func main() {
	engine := v0.RehearsalEngine{}

	fmt.Println(entity.GeneralPulicLicenseAbstruct)

	if en, err := load.Load(); err != nil {
		log.Fatal(err)
	} else {
		switch en.Command {
		case "run":
			if err := engine.AssignConfig(en); err != nil {
				log.Fatal(err)
			} else if err := engine.Run(); err != nil {
				log.Fatal(err)
			}
		case "about":
			fmt.Println(entity.AboutInfomation())
		}
	}
	engine.Finalize()
}
