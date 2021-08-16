// tasks/infrastructure.gateways/cui/cui_test.go
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

package cui_test

import (
	"os"

	cui_entity "github.com/rehearsal-open/rehearsal/tasks/infrastrcuture.gateways/cui.entity"
)

var (
	test_phase       = "test"
	test_cuitask1    = "cui_test1"
	test_cuitask2    = "cui_test2"
	only_output_task = cui_entity.Detail{
		Path: "python",
		Args: []string{
			"../../../test/python_to_python/src/python2.py",
		},
		Dir: "",
	}
	in_output_task = cui_entity.Detail{
		Path: "python",
		Args: []string{
			"../../../test/python_to_python/src/python1.py",
		},
		Dir: "",
	}
)

func init() {
	only_output_task.Dir, _ = os.Getwd()
}

func ErrorOccered(err error) {
	msg := err.Error()
	panic(msg)
}

// func TestEnv(t *testing.T) {
// 	cmd := exec.Command(only_output_task.Path, only_output_task.Args...)
// 	cmd.Dir = only_output_task.Dir

// 	if err := cmd.Run(); err != nil {
// 		ErrorOccered(err)
// 	}
// }

// func TestModule(t *testing.T) {

// 	if entity, err := rehearsal.AppendTask(test_phase, test_cuitask1, "cui", &only_output_task); err != nil {
// 		ErrorOccered(err)
// 	} else if val, err := cui.Make(entity); err != nil {
// 		ErrorOccered(err)
// 	} else if err := val.BeginTask(); err != nil {
// 		ErrorOccered(err)
// 	} else {
// 		val.WaitClosing()
// 		val.ReleaseResource()
// 	}

// 	return
// }

// func TestModule2(t *testing.T) {

// 	if entity1, err := rehearsal.AppendTask(test_phase, test_cuitask2, "cui", &only_output_task); err != nil {
// 		ErrorOccered(err)
// 	} else if entity2, err := rehearsal.AppendTask(test_phase, test_cuitask1, "cui", &in_output_task); err != nil {
// 		ErrorOccered(err)
// 	} else if task1, err := cui.Make(entity1); err != nil {
// 		ErrorOccered(err)
// 	} else if task2, err := cui.Make(entity2); err != nil {
// 		ErrorOccered(err)
// 	} else if reciever2, err := task2.Reciever(element.StdIn); err != nil {
// 		ErrorOccered(err)
// 	} else if err := task1.AppendReciever(element.StdOut, reciever2); err != nil {
// 		ErrorOccered(err)
// 	} else if err := task1.BeginTask(); err != nil {
// 		ErrorOccered(err)
// 	} else if err := task2.BeginTask(); err != nil {
// 		ErrorOccered(err)
// 	} else {
// 		task1.WaitClosing()
// 		task2.WaitClosing()
// 		task1.ReleaseResource()
// 		task2.ReleaseResource()
// 	}
// }
