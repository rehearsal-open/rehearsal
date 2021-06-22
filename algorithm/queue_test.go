package algorithm

import (
	"strconv"
	"testing"
)

func TestQueue(t *testing.T) {
	queue := SimpleQueueMaker(5)
	for i, j := 0, 0; i < 20; i, j = i+1, j+2 {
		if j < 20 {
			queue.Append(strconv.Itoa(j))
			queue.Append(strconv.Itoa(j + 1))
			// t.Log("queue: ", queue.data)
		}

		// var data string

		if data, err := queue.Get(); err != nil {
			t.Error(err)
		} else {
			num, _ := strconv.Atoi(data.(string))
			t.Log("queue[", i, "] = ", num+5, "(current index: ", queue.headIdx, ")")
			queue.Pop()
		}
	}
}
