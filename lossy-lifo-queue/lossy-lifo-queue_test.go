package lossylifoqueue

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

var comparator = func(a, b int) bool {
	return a == b
}

var strComparator = func(a, b string) bool {
	return a == b
}

var testData = []string{
	"/mnt/8tb-disk/photos",
	"/mnt/8tb-disk/videos",
	"/mnt/8tb-disk/documents",
	"/mnt/8tb-disk/photos",
	"/mnt/8tb-disk/documents/myresume.docx",
	"/mnt/8tb-disk/documents",
	"/mnt/8tb-disk/photos",
	"/mnt/8tb-disk/photos",
	"/mnt/8tb-disk/documents/myresume.docx",
}

func TestDuplicateElimination(t *testing.T) {
	q := NewLossyLifoQueue(5, strComparator)
	for _, s := range testData {
		bytes, _ := json.MarshalIndent(q, "", "  ")
		fmt.Printf("%s\n", string(bytes))
		json.Unmarshal(bytes, q)
		q.Add(s)
	}
	fmt.Printf("%v\n", q)
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

	peek_val, _ := queue.Peek()
	if peek_val != 4 {
		t.Fatalf("Peek expected %d, got %d", 4, peek_val)
	}

	pop_val, _ := queue.Pop()
	if pop_val != 4 {
		t.Fatalf("Pop expected %d, got %d", 4, pop_val)
	}

	expected = "[3, 2]"
	if queue.String() != expected {
		t.Fatalf("Expected %s, got %s", expected, fmt.Sprintf("%v", queue))
	}
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
	queueForSaving.Add(2)
	bytes, err := json.MarshalIndent(queueForSaving, "", " ")
	if err != nil {
		t.Fatalf("Unable to marshal queue: %v", err)
	}
	os.WriteFile(filename, bytes, 0644)

	var queueLoaded = NewLossyLifoQueue(3, comparator)
	readBytes, err := os.ReadFile(filename)
	log.Printf("Read %s\n", string(readBytes))
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
