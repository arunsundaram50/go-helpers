package lossylifoqueue

import (
	"container/list"
	"encoding/json"
	"fmt"
	"os"
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

// serialization part

type serializableQueue struct {
	Items   []interface{} `json:"items"`
	MaxSize int           `json:"maxSize"`
}

func (llq *LossyLifoQueue) Save(filename string) error {
	// Extract the items from the linked list into a slice.
	var items []interface{}
	for e := llq.data.Front(); e != nil; e = e.Next() {
		items = append(items, e.Value)
	}

	// Convert the list and maxSize into a serializable struct.
	toSave := serializableQueue{
		Items:   items,
		MaxSize: llq.maxSize,
	}

	// Serialize the struct to JSON.
	data, err := json.Marshal(toSave)
	if err != nil {
		return err
	}

	// Write the JSON data to the specified file.
	return os.WriteFile(filename, data, 0644)
}

func (llq *LossyLifoQueue) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	var loadedData serializableQueue
	err = json.Unmarshal(data, &loadedData)
	if err != nil {
		return err
	}

	// Clear the current queue.
	llq.data.Init()
	llq.lookup = make(map[interface{}]*list.Element)
	llq.maxSize = loadedData.MaxSize

	// Refill the queue using the loaded data.
	for _, item := range loadedData.Items {
		llq.Add(item)
	}

	return nil
}
