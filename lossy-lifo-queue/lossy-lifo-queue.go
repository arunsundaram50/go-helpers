package lossylifoqueue

import (
	"container/list"
	"fmt"
	"strings"
)

type LossyLifoQueue struct {
	data       *list.List
	maxSize    int
	comparator func(interface{}, interface{}) bool
	lookup     map[interface{}]*list.Element
}

func NewLossyLifoQueue(maxSize int, comparator func(interface{}, interface{}) bool) *LossyLifoQueue {
	return &LossyLifoQueue{
		data:       list.New(),
		maxSize:    maxSize,
		comparator: comparator,
		lookup:     make(map[interface{}]*list.Element),
	}
}

func (llq *LossyLifoQueue) Add(item interface{}) {
	// Check if item exists using the lookup map
	if elem, found := llq.lookup[item]; found {
		llq.data.Remove(elem)
	}

	// Add item to the end
	newElem := llq.data.PushBack(item)
	llq.lookup[item] = newElem

	// Check size constraints
	if llq.data.Len() > llq.maxSize {
		removedElem := llq.data.Front()
		llq.data.Remove(removedElem)
		delete(llq.lookup, removedElem.Value)
	}
}

func (llq *LossyLifoQueue) Pop() interface{} {
	if llq.data.Len() == 0 {
		return nil
	}
	lastElem := llq.data.Back()
	llq.data.Remove(lastElem)
	delete(llq.lookup, lastElem.Value)
	return lastElem.Value
}

func (llq *LossyLifoQueue) Peek() interface{} {
	if llq.data.Len() == 0 {
		return nil
	}
	return llq.data.Back().Value
}

func (llq *LossyLifoQueue) String() string {
	sb := strings.Builder{}
	sb.WriteString("[")
	firstItem := true
	for e := llq.data.Front(); e != nil; e = e.Next() {
		if firstItem {
			firstItem = false
		} else {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", e.Value))
	}
	sb.WriteString("]")
	return sb.String()
}
