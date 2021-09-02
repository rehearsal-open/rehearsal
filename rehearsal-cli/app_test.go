// rehearsal-cli/app_test.go
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
	"fmt"
	"testing"
)

// func BenchmarkWholeAlgorithm(t *testing.B) {

// 	t.StopTimer()

// 	for i, l := 0, t.N; i < l; i++ {
// 		parsemap := mapped.MappingType{
// 			"version": 0.202109,
// 			"phase": []interface{}{
// 				mapped.MappingType{
// 					"name": "phase_1",
// 					"task": []interface{}{
// 						mapped.MappingType{
// 							"name":      "python_1",
// 							"kind":      "cui",
// 							"wait-stop": true,
// 							"cmd":       "python",
// 							"args": []interface{}{
// 								"../test/py2py/01/python2.py",
// 							},
// 							"sendto": []interface{}{
// 								"phase_1::python_2",
// 								"phase_2::python_2",
// 							},
// 						},
// 						mapped.MappingType{
// 							"name":      "python_2",
// 							"kind":      "cui",
// 							"wait-stop": true,
// 							"cmd":       "python",
// 							"args": []interface{}{
// 								"../test/py2py/01/python1.py",
// 							},
// 							// "sendto": []
// 						},
// 					},
// 				},
// 				mapped.MappingType{
// 					"name": "phase_2",
// 					"task": []interface{}{
// 						// mapped.MappingType{
// 						// 	"name":      "python_1",
// 						// 	"kind":      "cui",
// 						// 	"wait-stop": true,
// 						// 	"cmd":       "python",
// 						// 	// "sendto": []
// 						// },
// 						mapped.MappingType{
// 							"name":      "python_2",
// 							"kind":      "cui",
// 							"wait-stop": true,
// 							"cmd":       "python",
// 							"args": []interface{}{
// 								"./test/py2py/01/python1.py",
// 							},
// 							"dir": "../",
// 							// "sendto": []
// 						},
// 					},
// 				},
// 			},
// 		}

// 		taskMaker := maker.Maker{}
// 		taskMaker.RegisterMaker("cui", &cui.MakeCollection)

// 		parser := mapped.Parser{
// 			DetailMaker: &taskMaker,
// 			Mapped:      parsemap,
// 		}

// 		en := engine.Rehearsal{}
// 		if err := en.Reset(&parser, taskMaker, nil); err != nil {
// 			panic(err)
// 		}

// 		t.StartTimer()

// 		if err := en.Execute(); err != nil {
// 			panic(err)
// 		}

// 		t.StopTimer()

// 	}

// 	t.Log("ok")
// }

func TestRun(t *testing.T) {
	if err := Run("../test/uart/01/rehearsal.yml"); err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
}
