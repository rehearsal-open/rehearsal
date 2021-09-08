// task/wrapper/splitter/splitter.go
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

package splitter

import (
	"io"
	"strings"

	"github.com/rehearsal-open/rehearsal/entities"
)

func (splitter *Splitter) Write(e *entities.Element, b []byte) error {

	splitter.cache += string(b)

	splitted := []string{splitter.cache}

	for _, splitStr := range splitter.SplitStr {
		splittedCache := splitted
		splitted = make([]string, 0)

		for _, val := range splittedCache {
			splitted = append(splitted, strings.Split(val, splitStr)...)
		}
	}

	for i, l := 0, len(splitted)-1; i < l; i++ {
		splitter.writer.Write([]byte(splitter.Prefix + splitted[i] + splitter.Suffix))
	}

	splitter.cache = splitted[len(splitted)-1]
	return nil
}

func (splitter *Splitter) Close() {
	splitter.writer.Write([]byte(splitter.cache))
}

func (splitter *Splitter) OutputTo(outto io.Writer) {
	splitter.writer = outto
}
