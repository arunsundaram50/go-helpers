package lossylifoqueue

import (
	"container/list"
	"encoding/json"
	"fmt"
	"strings"
)

type LossyLifoQueue struct {
	Data       *list.List
	MaxSize    int
	comparator func(interface{}, interface{}) bool
	lookup     map[interface{}]*list.Element
}

func NewLossyLifoQueue(maxSize int, comparator func(interface{}, interface{}) bool) *LossyLifoQueue {
	return &LossyLifoQueue{
		Data:       list.New(),
		MaxSize:    maxSize,
		comparator: comparator,
		lookup:     make(map[interface{}]*list.Element),
	}
}

func (llq *LossyLifoQueue) Add(item interface{}) {
	// Check if item exists using the lookup map
	if elem, found := llq.lookup[item]; found {
		llq.Data.Remove(elem)
	}

	// Add item to the end
	newElem := llq.Data.PushBack(item)
	llq.lookup[item] = newElem

	// Check size constraints
	if llq.Data.Len() > llq.MaxSize {
		removedElem := llq.Data.Front()
		llq.Data.Remove(removedElem)
		delete(llq.lookup, removedElem.Value)
	}
}

func (llq *LossyLifoQueue) Pop() interface{} {
	if llq.Data.Len() == 0 {
		return nil
	}
	lastElem := llq.Data.Back()
	llq.Data.Remove(lastElem)
	delete(llq.lookup, lastElem.Value)
	return lastElem.Value
}

func (llq *LossyLifoQueue) Peek() interface{} {
	if llq.Data.Len() == 0 {
		return nil
	}
	return llq.Data.Back().Value
}

func (llq *LossyLifoQueue) String() string {
	sb := strings.Builder{}
	sb.WriteString("[")
	firstItem := true
	for e := llq.Data.Front(); e != nil; e = e.Next() {
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

func (llq *LossyLifoQueue) GetAll() []interface{} {
	var items []interface{}
	for e := llq.Data.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value)
	}
	return items
}

func (llq *LossyLifoQueue) MarshalJSON() ([]byte, error) {
	// Extract items from the linked list into a slice for easier marshaling.
	var items []interface{}
	for e := llq.Data.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value)
	}

	// Create an auxiliary struct that represents the data we want to marshal.
	aux := struct {
		Items   []interface{} `json:"items"`
		MaxSize int           `json:"maxSize"`
	}{
		Items:   items,
		MaxSize: llq.MaxSize,
	}

	return json.Marshal(aux)
}

func (llq *LossyLifoQueue) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Items   []interface{} `json:"items"`
		MaxSize int           `json:"maxSize"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	llq.Data = list.New()
	for _, item := range aux.Items {
		llq.Data.PushBack(item)
	}
	llq.MaxSize = aux.MaxSize

	return nil
}
