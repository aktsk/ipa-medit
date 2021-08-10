package memory

import (
	"io/ioutil"
	"reflect"
	"testing"
)

func TestGetWritableAddrRanges(t *testing.T) {
	vmmapResult, _ := ioutil.ReadFile("testdata/vmmap_result")
	addrRanges, _ := GetWritableAddrRanges(vmmapResult)

	actualLens := len(addrRanges)
	expectedLens := 1937
	if actualLens != expectedLens {
		t.Errorf("got length: %v\nexpected length: %v", actualLens, expectedLens)
	}

	actualLastRange := addrRanges[len(addrRanges)-1]
	expectedLastRange := [2]int{105553518919680, 105553653137408}
	if !reflect.DeepEqual(actualLastRange, expectedLastRange) {
		t.Errorf("got last address range: %v\nexpected last address range: %v", actualLastRange, expectedLastRange)
	}
}

func TestFindDataInSplittedMemory(t *testing.T) {
	memory := []byte{0x10, 0x11, 0x12, 0x10, 0x10, 0x11, 0x12, 0x11, 0x10, 0x11, 0x12, 0x12}
	searchBytes := []byte{0x10, 0x11, 0x12}
	actual := []int{}
	findDataInSplittedMemory(&memory, searchBytes, len(searchBytes), 0x100, 0, &actual)
	expected := []int{0x100, 0x104, 0x108}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got addr slice: %v\nexpected addr slice: %v", actual, expected)
	}
}

func TestFindEmptyInSplittedMemory(t *testing.T) {
	memory := []byte{0x10}
	searchBytes := []byte{0xAA, 0xBB, 0xCC}
	actual := []int{}
	findDataInSplittedMemory(&memory, searchBytes, len(searchBytes), 0x100, 0, &actual)
	expected := []int{}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got addr slice: %v\nexpected addr slice: %v", actual, expected)
	}
}
