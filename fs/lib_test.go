package fs

import (
	"os"
	"testing"
	"time"
)

func TesFileModifiedTime(t *testing.T) {
	filename := "/tmp/a.txt"
	if err := Touch(filename); err != nil {
		t.Fatal(err)
	}

	if d, err := GetDurationSinceModified("/tmp/a.txt"); err != nil {
		t.Fatal(err)
	} else {
		if d > time.Minute {
			t.Fatalf("expected < a minute, but got %v", d)
		}
	}
}

func TestIsStale(t *testing.T) {
	srcFilename := "/tmp/a.txt"
	derFilename := "/tmp/a-derived.txt"

	os.WriteFile(srcFilename, []byte("Hello, world"), 0644)
	time.Sleep(100 * time.Millisecond)

	readBytes, _ := os.ReadFile(srcFilename)
	os.WriteFile(derFilename, readBytes, 0644)

	isStale, _ := IsStale(srcFilename, derFilename)
	if isStale {
		t.Fatalf("File %s should not be stale compared to %s", derFilename, srcFilename)
	}
}
