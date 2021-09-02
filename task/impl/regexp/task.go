// task/impl/regexp/task.go
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

package regexp

import (
	"bytes"
	"io"

	"github.com/pkg/errors"
	"github.com/rehearsal-open/rehearsal/entities/enum/task_element"
	"github.com/rehearsal-open/rehearsal/task/based"
	"github.com/rehearsal-open/rehearsal/task/buffer"
)

func (matching *__task) IsSupporting(elem task_element.Enum) bool {
	return [task_element.Len]bool{
		true, true, false,
	}[elem]
}

func (matching *__task) ExecuteMain(args based.MainFuncArguments) error {

	var stdOut io.Writer

	if stdout, err := args.Writer(task_element.StdOut); err != nil {
		return errors.WithStack(err)
	} else {
		stdOut = stdout
	}

	cache := bytes.NewBufferString("")

	callback := [task_element.Len]based.ImplCallback{nil}
	callback[task_element.StdIn] = based.MakeImplCallback(func(recieved *buffer.Packet) {
		defer recieved.Close()

		io.Copy(cache, recieved)
		buffer := cache.String()

		matches := matching.Regexp.FindAllStringSubmatchIndex(buffer, -1)
		read := 0

		for i := range matches {
			results := matching.Regexp.ExpandString([]byte{}, matching.detail.TemplateRegexpr, buffer, matches[i])
			stdOut.Write(results)
			if matches[i][1] > read {
				read = matches[i][1]
			}
		}

		if read > 0 {
			cache = bytes.NewBufferString(buffer[read:])
		}

	}, func() {
		if matching.detail.FinalForcePush && cache.Len() > 0 {
			stdOut.Write(cache.Bytes())
		}
	})

	go func() {
		args.ListenStart(callback)
		<-matching.close
		args.Close(nil)
	}()

	return nil
}

func (matching *__task) StopMain() {
	close(matching.close)
}
