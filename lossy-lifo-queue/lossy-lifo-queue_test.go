package lossylifoqueue

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

var comparator = func(a, b interface{}) bool {
	return a.(int) == b.(int)
}

func TestIntOperations(t *testing.T) {
	queue := NewLossyLifoQueue(3, comparator)
	queue.Add(1)
	queue.Add(2)
	queue.Add(3)
	expected := "[1, 2, 3]"
	if queue.String() != expected {
		t.Fatalf("Expected %s, got %s", expected, fmt.Sprintf("%v", queue))
	}

	queue.Add(2) // remove existing 2 and add to end
	expected = "[1, 3, 2]"
	if queue.String() != expected {
		t.Fatalf("Expected %s, got %s", expected, fmt.Sprintf("%v", queue))
	}

	queue.Add(4)
	expected = "[3, 2, 4]"
	if queue.String() != expected {
		t.Fatalf("Expected %s, got %s", expected, fmt.Sprintf("%v", queue))
	}

	peek_val := queue.Peek()
	if peek_val != 4 {
		t.Fatalf("Peek expected %d, got %d", 4, peek_val)
	}

	pop_val := queue.Pop()
	if pop_val != 4 {
		t.Fatalf("Pop expected %d, got %d", 4, pop_val)
	}

	expected = "[3, 2]"
	if queue.String() != expected {
		t.Fatalf("Expected %s, got %s", expected, fmt.Sprintf("%v", queue))
	}
}

func strComparator(a, b interface{}) bool {
	return a.(string) == b.(string)
}

func TestStringOperations(t *testing.T) {
	lastDir := NewLossyLifoQueue(10, strComparator)
	lastDir.Add("hello")
	lastDir.Add("World")
	fmt.Println(lastDir)
}

func TestSaveLoad(t *testing.T) {
	filename := "/tmp/x.json"
	queueForSaving := NewLossyLifoQueue(3, comparator)
	queueForSaving.Add(3)
	queueForSaving.Add(2)
	bytes, err := json.MarshalIndent(queueForSaving, "", " ")
	if err != nil {
		t.Fatalf("Unable to marshal queue: %v", err)
	}
	os.WriteFile(filename, bytes, 0644)

	var queueLoaded = NewLossyLifoQueue(3, comparator)
	readBytes, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Unable to read %s: %v", filename, err)
	}

	if err := json.Unmarshal(readBytes, &queueLoaded); err != nil {
		t.Fatalf("Unable to unmarshal %s: %v", filename, err)
	}

	expected := "[3, 2]"
	if queueLoaded.String() != expected {
		t.Fatalf("Expected %s, got %s", expected, fmt.Sprintf("%v", queueLoaded))
	}

}
