package memory

import (
	"reflect"
	"testing"
	"io/ioutil"
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