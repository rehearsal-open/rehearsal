// tasks/infrastructure.gateways/cui/cui_benchmark_test.go
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
	"testing"

	"github.com/rehearsal-open/rehearsal/entities"
	"github.com/rehearsal-open/rehearsal/entities/element"
	"github.com/rehearsal-open/rehearsal/tasks/infrastructure.gateways/cui"
)

func BenchmarkModule2(b *testing.B) {

	b.Log("loop times: ", b.N)

	for i, l := 0, b.N; i < l; i++ {

		rehearsal := entities.NewRehearsal()
		rehearsal.AppendPhase(test_phase)

		if entity1, err := rehearsal.AppendTask(test_phase, test_cuitask2, "cui", &only_output_task); err != nil {
			ErrorOccered(err)
		} else if entity2, err := rehearsal.AppendTask(test_phase, test_cuitask1, "cui", &in_output_task); err != nil {
			ErrorOccered(err)
		} else if task1, err := cui.Make(entity1); err != nil {
			ErrorOccered(err)
		} else if task2, err := cui.Make(entity2); err != nil {
			ErrorOccered(err)
		} else if reciever2, err := task2.Reciever(element.StdIn); err != nil {
			ErrorOccered(err)
		} else if err := task1.AppendReciever(element.StdOut, reciever2); err != nil {
			ErrorOccered(err)
		} else if err := task1.BeginTask(); err != nil {
			ErrorOccered(err)
		} else if err := task2.BeginTask(); err != nil {
			ErrorOccered(err)
		} else {
			task1.WaitClosing()
			task2.WaitClosing()
			task1.ReleaseResource()
			task2.ReleaseResource()
		}
	}
}
