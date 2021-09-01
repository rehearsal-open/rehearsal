// .for-dev/makefile/makefile.go
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
	"os"
	"path/filepath"
	"strings"
	"time"
)

const initFile string = `// %s
// Copyright (C) %d %s

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

package %s


`

func main() {

	// get command environment
	wd, _ := os.Getwd()
	if len(os.Args) < 3 {
		log.Fatalln("command has a few arguments....",
			"this cli use: make-gofile <relative file path> <author name> (<package name>)")
	} else if len(os.Args) > 4 {
		log.Fatalln("command has much arguments....",
			"this cli use: make-gofile <relative file path> <author name> (<package name>)")
	}
	relPath := os.Args[1]
	author := os.Args[2]
	path := filepath.Join(wd, relPath)
	dir := filepath.Dir(path)
	var pkgName string
	if len(os.Args) == 4 {
		pkgName = os.Args[3]
	} else {
		pkgName = strings.Split(filepath.Base(dir), ".")[0]
	}

	// check directory and filename
	if f, err := os.Stat(filepath.Join(wd, "/.git")); os.IsNotExist(err) || !f.IsDir() {
		log.Fatalln("you must use in root directory of repository....",
			fmt.Sprint("current directory: ", wd))
	} else if f, err := os.Stat(dir); os.IsNotExist(err) || !f.IsDir() {
		log.Fatalln("you must select already made directory; make directory....\n",
			fmt.Sprint("your select file: ", dir))
	} else if f, err := os.Stat(path); !(os.IsNotExist(err) || f.IsDir()) {
		log.Fatalln("your select file has already existed....")
	} else if filepath.Ext(path) != ".go" {
		log.Fatalln("this script makes only .go files....")
	}

	// make & write file
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	output := fmt.Sprintf(initFile, relPath, time.Now().UTC().Year(), author, pkgName)
	file.Write(([]byte)(output))

	print(`create files: success!
	File path:       ` + relPath + `
	  (package name: ` + pkgName + `)
	Author Name:     ` + author)
}
