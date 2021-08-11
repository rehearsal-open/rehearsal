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

	// Command Documentation
	// make-gofile [root directory path] [adding *.go file's relative path] [author name] ([module name])

	// get command environment
	wd, _ := os.Getwd()

	// check number of arguments
	if len(os.Args) < 4 {
		log.Fatalln("command has a few arguments....",
			"this cli use: make-gofile <relative file path> <author name> (<package name>)")
	} else if len(os.Args) > 5 {
		log.Fatalln("command has much arguments....",
			"this cli use: make-gofile <relative file path> <author name> (<package name>)")
	}

	// variable definates
	var (
		absPath    string // absolute path
		absPathDir string
		relPath    string // relative path, base directory is root directory
		pkgName    string // package name
		rootDir    string // root of repository's absolute path
		author     string = os.Args[3]
	)

	// making absolute path
	if abs, err := filepath.Abs(wd); err != nil {
		log.Fatalln(err.Error)
	} else {
		absPath = filepath.Join(abs, os.Args[2])
		absPathDir = filepath.Dir(absPath)
	}

	// making root directory path
	if root, err := filepath.Abs(os.Args[1]); err != nil {
		log.Fatalln(err.Error())
	} else if f, err := os.Stat(root); os.IsNotExist(err) || !f.IsDir() {
		log.Fatalln("you must use in root directory of repository....",
			fmt.Sprint("current directory: ", wd))
	} else {
		rootDir = root
	}

	// making relative path
	if f, err := os.Stat(filepath.Join(rootDir, "/.git")); os.IsNotExist(err) || !f.IsDir() {
		log.Fatalln("you must use in root directory of repository....",
			fmt.Sprint("current directory: ", wd))
	} else if rel, err := filepath.Rel(rootDir, absPath); err != nil {
		log.Fatalln(err.Error())
	} else {
		relPath = strings.ReplaceAll(rel, "\\", "/")
	}

	// making package name
	if len(os.Args) == 5 {
		pkgName = os.Args[4]
	} else {
		pkgName = strings.Split(filepath.Base(absPathDir), ".")[0]
	}

	// check directory and filename
	if f, err := os.Stat(absPathDir); os.IsNotExist(err) || !f.IsDir() {
		log.Fatalln("you must select already made directory; make directory....\n",
			fmt.Sprint("your select file: ", absPathDir))
	} else if f, err := os.Stat(absPath); !(os.IsNotExist(err) || f.IsDir()) {
		log.Fatalln("your select file has already existed....")
	} else if filepath.Ext(absPath) != ".go" {
		log.Fatalln("this script makes only .go files....")
	}

	// make & write file
	file, err := os.Create(absPath)
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
