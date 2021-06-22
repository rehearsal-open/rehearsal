package algorithm

import (
	"github.com/pkg/errors"
)

// Queue algorithm object
type SimpleQueue struct {
	data    []interface{}
	headIdx uint64
	tailIdx uint64
	length  uint64
}

// Constractor of algorithm.SingleQueue.
func SimpleQueueMaker(margin uint64) SimpleQueue {

	queue := SimpleQueue{
		data:    make([]interface{}, margin),
		headIdx: 0,
		tailIdx: 0,
		length:  0,
	}

	return queue
}

// Return length of queue's internal array (internal method).
func (queue *SimpleQueue) ability() uint64 {
	return uint64(len(queue.data))
}

// Return result of increment index (internal method).
func (queue *SimpleQueue) increment(idx uint64) uint64 {
	idx++
	if idx == queue.ability() {
		idx = uint64(0)
	}
	return idx
}

// Return result of addition index (internal method).
func (queue *SimpleQueue) add(idx uint64, addVal uint64) uint64 {
	ability := queue.ability()
	remain := ability - idx
	if addVal < remain {
		idx += addVal
	} else {
		idx = (addVal - remain) % ability
	}
	return idx
}

// Push back the object into queue.
func (queue *SimpleQueue) Append(object interface{}) {

	ability := queue.ability()
	queue.length++

	if queue.length > ability {

		new := make([]interface{}, ability<<1)
		for iSrc, iDst := queue.headIdx, uint64(0); iDst+1 < queue.length; iSrc, iDst = queue.increment(iSrc), iDst+1 {
			new[iDst] = queue.data[iSrc]
		}

		queue.headIdx = 0
		queue.tailIdx = queue.length
		new[queue.length-1] = object
		queue.data = new

	} else {

		queue.data[queue.tailIdx] = object
		queue.tailIdx = queue.increment(queue.tailIdx)

	}
}

// Return the front value; this func doesn't pop (i.e: don't remove front value from queue).
func (queue *SimpleQueue) Get() (interface{}, error) {
	if queue.length < 1 {
		return nil, errors.New("Queue is empty")
	}
	return queue.data[queue.headIdx], nil
}

// Return the front value and remove front value from queue.
func (queue *SimpleQueue) GetPop() (result interface{}, err error) {
	result, err = queue.Get()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	queue.Pop()
	return result, err
}

// Remove front value from queue.
func (queue *SimpleQueue) Pop() {
	if queue.length > 0 {
		queue.length--
		queue.headIdx = queue.increment(queue.headIdx)
	}
}

// Return length of queue, it is on the surface.
func (queue *SimpleQueue) Length() uint64 {
	return queue.length
}
