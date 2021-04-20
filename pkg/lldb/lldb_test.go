package lldb

import (
	"os"
	"testing"
)

func TestFileExistsWhenExists(t *testing.T) {
	exePath, _ := os.Executable()
	actual := fileExists(exePath)
	expected := true
	if actual != expected {
		t.Errorf("got: %v\nexpected: %v", actual, expected)
	}
}

func TestFileExistsWhenNotExists(t *testing.T) {
	actual := fileExists("./a.py")
	expected := false
	if actual != expected {
		t.Errorf("got: %v\nexpected: %v", actual, expected)
	}
}
