package lossylifoqueue

import (
	"container/list"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type LossyLifoQueue[T comparable] struct {
	zeroValue  T
	Data       *list.List
	MaxSize    int
	comparator func(T, T) bool
	lookup     map[T]*list.Element
}

func NewLossyLifoQueue[T comparable](maxSize int, comparator func(T, T) bool) *LossyLifoQueue[T] {
	return &LossyLifoQueue[T]{
		Data:       list.New(),
		MaxSize:    maxSize,
		comparator: comparator,
		lookup:     make(map[T]*list.Element),
	}
}

func (llq *LossyLifoQueue[T]) Add(item T) {
	// Check if item exists using the lookup map
	if elem, found := llq.lookup[item]; found {
		llq.Data.Remove(elem)
		delete(llq.lookup, item) // GOTCHA deleting from the map helps keep the list and map in sync so that the code below does not cause duplicates
	}

	// Add item to the end
	newElem := llq.Data.PushBack(item)
	llq.lookup[item] = newElem

	// Check size constraints
	if llq.Data.Len() > llq.MaxSize {
		removedElem := llq.Data.Front()
		llq.Data.Remove(removedElem)
		if removedValue, ok := removedElem.Value.(T); ok {
			delete(llq.lookup, removedValue)
		}
	}
}

func (llq *LossyLifoQueue[T]) Pop() (T, bool) {
	if llq.Data.Len() == 0 {
		return llq.zeroValue, false
	}
	lastElem := llq.Data.Back()
	llq.Data.Remove(lastElem)
	if lastValue, ok := lastElem.Value.(T); ok {
		delete(llq.lookup, lastValue)
	}

	if lastValue, ok := lastElem.Value.(T); ok {
		delete(llq.lookup, lastValue)
		return lastValue, true
	}
	return llq.zeroValue, false
}

func (llq *LossyLifoQueue[T]) Peek() (T, bool) {
	if llq.Data.Len() == 0 {
		return llq.zeroValue, false
	}
	lastElem := llq.Data.Back()

	// Assert the type of lastElem.Value to T
	if lastValue, ok := lastElem.Value.(T); ok {
		return lastValue, true
	}

	return llq.zeroValue, false
}

func (llq *LossyLifoQueue[T]) String() string {
	sb := strings.Builder{}
	sb.WriteString("[\n")
	firstItem := true
	for e := llq.Data.Front(); e != nil; e = e.Next() {
		if firstItem {
			firstItem = false
		} else {
			sb.WriteString(",\n")
		}
		sb.WriteString(fmt.Sprintf("  %v", e.Value))
	}
	sb.WriteString("\n]")
	return sb.String()
}

func (llq *LossyLifoQueue[T]) GetAll() []T {
	var items []T
	for e := llq.Data.Front(); e != nil; e = e.Next() {
		if v, ok := e.Value.(T); ok {
			items = append(items, v)
		}
	}
	return items
}

// Showoff leaving the MarshalJSON and UnmarshalJSON to use interface{} is a practical and common approach, especially when dealing with generics
// This approach simplifies the serialization process and works well in many cases
func (llq *LossyLifoQueue[T]) MarshalJSON() ([]byte, error) {
	// Extract items from the linked list into a slice for easier marshaling.
	var items []interface{}
	for e := llq.Data.Front(); e != nil; e = e.Next() {
		// Append
		// items = append(items, e.Value)
		// Prepend
		items = append([]interface{}{e.Value}, items...)
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

func (llq *LossyLifoQueue[T]) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Items   []interface{} `json:"items"`
		MaxSize int           `json:"maxSize"`
	}{}

	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	llq.Data = list.New()
	llq.lookup = make(map[T]*list.Element) // GOTCHA Not reinitializing the lookup map would make duplicate elimination impossible!

	for i := len(aux.Items) - 1; i >= 0; i-- {
		item := aux.Items[i]
		typedItem, ok := item.(T)
		if !ok {
			log.Printf("incorrect type for %v\n", typedItem)
			continue
		}

		newElem := llq.Data.PushBack(typedItem)
		llq.lookup[typedItem] = newElem
	}
	llq.MaxSize = aux.MaxSize

	return nil
}
