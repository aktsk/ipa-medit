package memory

import (
	"bufio"
	"strconv"
	"bytes"
)

func GetWritableAddrRanges(vmmapResult []byte) ([][2]int, error) {
	addrRanges := [][2]int{}
	//ignorePaths := []string{"/vendor/lib64/", "/system/lib64/", "/system/bin/", "/system/framework/", "/data/dalvik-cache/"}
	scanner := bufio.NewScanner(bytes.NewReader(vmmapResult))
	writable := false
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.HasPrefix(line, []byte("==== Writable regions for process")) {
			writable = true
			continue
		} else if len(line) <= 1 {
			writable = false
		}

		if writable && !bytes.HasPrefix(line, []byte("REGION TYPE")) {
			meminfo := bytes.Fields(line)
			region_type := meminfo[0]
			i := 0
			for i = 0; i < 3; i++ {
				if (bytes.Index(line, meminfo[i+1]) - bytes.Index(line, meminfo[i]) - len(meminfo[i])) == 1 {
					region_type = append(region_type, meminfo[i+1]...)
				} else {
					break
				}
			}
			addrRange := meminfo[i+1]
			ignoreFlag := false
			if !ignoreFlag {
				addrs := bytes.Split(addrRange, []byte("-"))
				beginAddr, _ := strconv.ParseInt(string(addrs[0]), 16, 64)
				endAddr, _ := strconv.ParseInt(string(addrs[1]), 16, 64)
				addrRanges = append(addrRanges, [2]int{int(beginAddr), int(endAddr)})
			}
		}
	}
	return addrRanges, nil
}