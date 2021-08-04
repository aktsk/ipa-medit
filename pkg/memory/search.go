package memory

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/aktsk/ipa-medit/pkg/converter"
)

func GetWritableAddrRanges(vmmapResult []byte) ([][2]int, error) {
	addrRanges := [][2]int{}
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
			addrs := bytes.Split(addrRange, []byte("-"))
			beginAddr, _ := strconv.ParseInt(string(addrs[0]), 16, 64)
			endAddr, _ := strconv.ParseInt(string(addrs[1]), 16, 64)
			addrRanges = append(addrRanges, [2]int{int(beginAddr), int(endAddr)})
		}
	}
	return addrRanges, nil
}

var splitSize = 0x5000000
var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, splitSize)
	},
}

type Err struct {
	err error
}

func (e *Err) Error() string {
	return fmt.Sprint(e.err)
}

type ParseErr struct {
	*Err
}

type TooManyErr struct {
	*Err
}

func FindDataInAddrRanges(pid int, targetBytes []byte, addrRanges [][2]int) ([]int, error) {
	foundAddrs := []int{}
	searchLength := len(targetBytes)
	for _, s := range addrRanges {
		beginAddr := s[0]
		endAddr := s[1]
		memSize := endAddr - beginAddr
		for i := 0; i < (memSize/splitSize)+1; i++ {
			// target memory is too big to read all of it, so split it and then search in memory
			splitIndex := (i + 1) * splitSize
			splittedBeginAddr := beginAddr + i*splitSize
			splittedEndAddr := endAddr
			if splitIndex < memSize {
				splittedEndAddr = beginAddr + splitIndex
			}
			b := bufferPool.Get().([]byte)[:(splittedEndAddr - splittedBeginAddr)]
			task := GetTaskForPid(pid)
			if err := ReadMemory(task, b, splittedBeginAddr, splittedEndAddr); err == nil {
				findDataInSplittedMemory(&b, targetBytes, searchLength, splittedBeginAddr, 0, &foundAddrs)
				bufferPool.Put(b)
				if len(foundAddrs) > 500000 {
					fmt.Println("Too many addresses with target data found...")
					return foundAddrs, TooManyErr{&Err{errors.New("Error: Too many addresses")}}
				}
			} else {
				fmt.Printf("0x%x: %s\n", beginAddr, err)
			}
		}
	}
	return foundAddrs, nil
}

func findDataInSplittedMemory(memory *[]byte, targetBytes []byte, searchLength int, beginAddr int, offset int, results *[]int) {
	// use Rabin-Karp string search algorithm in bytes.Index
	index := bytes.Index((*memory)[offset:], targetBytes)
	if index == -1 {
		return
	} else {
		resultAddr := beginAddr + index + offset
		*results = append(*results, resultAddr)
		offset += index + searchLength
		findDataInSplittedMemory(memory, targetBytes, searchLength, beginAddr, offset, results)
	}
}

func FindString(pid int, targetVal string, addrRanges [][2]int) ([]int, error) {
	fmt.Println("Search UTF-8 String...")
	targetBytes, _ := converter.StringToBytes(targetVal)
	fmt.Printf("Target Value: %s(%v)\n", targetVal, targetBytes)
	foundAddrs, err := FindDataInAddrRanges(pid, targetBytes, addrRanges)
	fmt.Printf("Found: %d!!\n", len(foundAddrs))
	if len(foundAddrs) < 10 {
		for _, v := range foundAddrs {
			fmt.Printf("Address: 0x%x\n", v)
		}
	}
	return foundAddrs, err
}

func FindWord(pid int, targetVal string, addrRanges [][2]int) ([]int, error) {
	fmt.Println("Search Word...")
	targetBytes, err := converter.WordToBytes(targetVal)
	if err != nil {
		fmt.Printf("parsing %s: value out of range\n", targetVal)
		return nil, ParseErr{&Err{errors.New("Error: value out of range")}}
	}
	fmt.Printf("Target Value: %s(%v)\n", targetVal, targetBytes)
	foundAddrs, err := FindDataInAddrRanges(pid, targetBytes, addrRanges)
	fmt.Printf("Found: %d!!\n", len(foundAddrs))
	if len(foundAddrs) < 10 {
		for _, v := range foundAddrs {
			fmt.Printf("Address: 0x%x\n", v)
		}
	}
	return foundAddrs, err
}

func FindDword(pid int, targetVal string, addrRanges [][2]int) ([]int, error) {
	fmt.Println("Search Double Word...")
	targetBytes, err := converter.DwordToBytes(targetVal)
	if err != nil {
		fmt.Printf("parsing %s: value out of range\n", targetVal)
		return nil, ParseErr{&Err{errors.New("Error: value out of range")}}
	}
	fmt.Printf("Target Value: %s(%v)\n", targetVal, targetBytes)
	foundAddrs, err := FindDataInAddrRanges(pid, targetBytes, addrRanges)
	fmt.Printf("Found: %d!!\n", len(foundAddrs))
	if len(foundAddrs) < 10 {
		for _, v := range foundAddrs {
			fmt.Printf("Address: 0x%x\n", v)
		}
	}
	return foundAddrs, err
}

func FindQword(pid int, targetVal string, addrRanges [][2]int) ([]int, error) {
	fmt.Println("Search Quad Word...")
	targetBytes, err := converter.QwordToBytes(targetVal)
	if err != nil {
		fmt.Printf("parsing %s: value out of range\n", targetVal)
		return nil, ParseErr{&Err{errors.New("Error: value out of range")}}
	}
	fmt.Printf("Target Value: %s(%v)\n", targetVal, targetBytes)
	foundAddrs, err := FindDataInAddrRanges(pid, targetBytes, addrRanges)
	fmt.Printf("Found: %d!!\n", len(foundAddrs))
	if len(foundAddrs) < 10 {
		for _, v := range foundAddrs {
			fmt.Printf("Address: 0x%x\n", v)
		}
	}
	return foundAddrs, err
}
