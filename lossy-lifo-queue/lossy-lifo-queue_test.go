package lossylifoqueue

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	comparator := func(a, b interface{}) bool {
		return a.(int) == b.(int)
	}

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
