// rehearsal-cli/main.go
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

import (
	"flag"
	"fmt"
)

var (
	cmd string
)

func init() {
	init_run()
	flag.Usage = PutHelpDefault
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("This program is required subcommands!")
		PutHelpDefault()
	}

	switch flag.Arg(0) {
	case "run":
		if flag.NArg() < 2 {
			fmt.Println("'rehearsal-cli run' command is require config file's path")
			fmt.Println(helpRun)
			return
		} else if sub := flag.Arg(1); sub == "help" {
			fmt.Println(helpRun)
			return
		} else if err := Run(sub); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}
