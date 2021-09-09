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

	splitter.lock.Lock()
	defer splitter.lock.Unlock()

	if splitter.cache == nil {
		splitter.cache = make([]byte, 0, 1024)
		splitter.buffer = make([]byte, 0, 1024)
		splitter.buffer = append(splitter.buffer, []byte(splitter.Prefix)...)
	}

	splitter.cache = append(splitter.cache, b...)
	begins := 0

	for isContinue := true; isContinue && begins < len(splitter.cache); {
		isContinue = false
		for _, splitStr := range splitter.SplitStr {
			if ends := strings.Index(string(splitter.cache[begins:]), splitStr); ends > -1 {

				// when splitting string is found
				isContinue = true
				ends += begins

				// append bytes
				splitter.buffer = append(splitter.buffer, splitter.cache[begins:ends]...)
				splitter.buffer = append(splitter.buffer, []byte(splitter.Suffix)...)
				splitter.writer.Write(splitter.buffer)

				// reset bytes
				splitter.buffer = splitter.buffer[:len(splitter.Prefix)]
				begins += ends + len(splitStr)
				break
			}
		}
	}

	// cache set
	if begins > 0 {
		if begins < len(splitter.cache) {
			copy(splitter.cache, splitter.cache[begins:])
			splitter.cache = splitter.cache[:len(splitter.cache)-begins]
		} else {
			splitter.cache = splitter.cache[:0]
		}

	}

	return nil
}

func (splitter *Splitter) Close() {
	splitter.lock.Lock()
	defer splitter.lock.Unlock()
	splitter.writer.Write([]byte(splitter.cache))
}

func (splitter *Splitter) OutputTo(outto io.Writer) {
	splitter.writer = outto
}
