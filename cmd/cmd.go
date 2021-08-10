package cmd

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os/exec"
	"strconv"

	sys "golang.org/x/sys/unix"

	"github.com/aktsk/ipa-medit/pkg/converter"
	"github.com/aktsk/ipa-medit/pkg/memory"
)

type Found struct {
	addrs     []int
	converter func(string) ([]byte, error)
	dataType  string
}

var isAttached = false

func Find(pid string, targetVal string, dataType string) ([]Found, error) {
	founds := []Found{}
	result, err := exec.Command("vmmap", "--wide", pid).Output()
	if err != nil {
		return nil, err
	}
	addrRanges, _ := memory.GetWritableAddrRanges(result)

	var intPid int
	if intPid, err = strconv.Atoi(pid); err != nil {
		return nil, err
	}

	if dataType == "all" {
		// search string
		foundAddrs, err := memory.FindString(intPid, targetVal, addrRanges)
		if err == nil && len(foundAddrs) > 0 {
			founds = append(founds, Found{
				addrs:     foundAddrs,
				converter: converter.StringToBytes,
				dataType:  "UTF-8 string",
			})
		} else if _, ok := err.(memory.TooManyErr); ok {
			return founds, err
		}
		fmt.Println("------------------------")

		// search int
		foundAddrs, err = memory.FindWord(intPid, targetVal, addrRanges)
		if err == nil {
			if len(foundAddrs) > 0 {
				founds = append(founds, Found{
					addrs:     foundAddrs,
					converter: converter.WordToBytes,
					dataType:  "word",
				})
			}
			return founds, nil
		} else if _, ok := err.(memory.TooManyErr); ok {
			return founds, err
		}
		fmt.Println("------------------------")
		foundAddrs, err = memory.FindDword(intPid, targetVal, addrRanges)
		if err == nil {
			if len(foundAddrs) > 0 {
				founds = append(founds, Found{
					addrs:     foundAddrs,
					converter: converter.DwordToBytes,
					dataType:  "dword",
				})
			}
			return founds, nil
		} else if _, ok := err.(memory.TooManyErr); ok {
			return founds, err
		}
		fmt.Println("------------------------")
		foundAddrs, err = memory.FindQword(intPid, targetVal, addrRanges)
		if err == nil {
			if len(foundAddrs) > 0 {
				founds = append(founds, Found{
					addrs:     foundAddrs,
					converter: converter.QwordToBytes,
					dataType:  "qword",
				})
			}
			return founds, nil
		} else if _, ok := err.(memory.TooManyErr); ok {
			return founds, err
		}

	} else if dataType == "string" {
		foundAddrs, _ := memory.FindString(intPid, targetVal, addrRanges)
		if err == nil {
			if len(foundAddrs) > 0 {
				founds = append(founds, Found{
					addrs:     foundAddrs,
					converter: converter.StringToBytes,
					dataType:  "UTF-8 string",
				})
			}
			return founds, nil
		} else if _, ok := err.(memory.TooManyErr); ok {
			return founds, err
		}

	} else if dataType == "word" {
		foundAddrs, err := memory.FindWord(intPid, targetVal, addrRanges)
		if err == nil {
			if len(foundAddrs) > 0 {
				founds = append(founds, Found{
					addrs:     foundAddrs,
					converter: converter.WordToBytes,
					dataType:  "word",
				})
			}
			return founds, nil
		} else if _, ok := err.(memory.TooManyErr); ok {
			return founds, err
		}

	} else if dataType == "dword" {
		foundAddrs, err := memory.FindDword(intPid, targetVal, addrRanges)
		if err == nil {
			if len(foundAddrs) > 0 {
				founds = append(founds, Found{
					addrs:     foundAddrs,
					converter: converter.DwordToBytes,
					dataType:  "dword",
				})
			}
			return founds, nil
		} else if _, ok := err.(memory.TooManyErr); ok {
			return founds, err
		}

	} else if dataType == "qword" {
		foundAddrs, err := memory.FindQword(intPid, targetVal, addrRanges)
		if err == nil {
			if len(foundAddrs) > 0 {
				founds = append(founds, Found{
					addrs:     foundAddrs,
					converter: converter.QwordToBytes,
					dataType:  "qword",
				})
			}
			return founds, nil
		} else if _, ok := err.(memory.TooManyErr); ok {
			return founds, err
		}
	}

	return nil, errors.New("Error: specified datatype does not exist")
}

func Filter(pid string, targetVal string, prevFounds []Found) ([]Found, error) {
	founds := []Found{}
	result, err := exec.Command("vmmap", "--wide", pid).Output()
	if err != nil {
		return nil, err
	}
	writableAddrRanges, err := memory.GetWritableAddrRanges(result)
	if err != nil {
		return nil, err
	}

	var intPid int
	if intPid, err = strconv.Atoi(pid); err != nil {
		return nil, err
	}
	addrRanges := [][2]int{}

	// check if previous result address exists in current memory map
	for i, prevFound := range prevFounds {
		targetBytes, _ := prevFound.converter(targetVal)
		targetLength := len(targetBytes)
		fmt.Printf("Check previous results of searching %s...\n", prevFound.dataType)
		fmt.Printf("Target Value: %s(%v)\n", targetVal, targetBytes)
		for _, prevAddr := range prevFound.addrs {
			for _, writable := range writableAddrRanges {
				if writable[0] < prevAddr && prevAddr < writable[1] {
					addrRanges = append(addrRanges, [2]int{prevAddr, prevAddr + targetLength})
				}
			}
		}
		foundAddrs, _ := memory.FindDataInAddrRanges(intPid, targetBytes, addrRanges)
		fmt.Printf("Found: %d!!\n", len(foundAddrs))
		if len(foundAddrs) < 10 {
			for _, v := range foundAddrs {
				fmt.Printf("Address: 0x%x\n", v)
			}
		}
		founds = append(founds, Found{
			addrs:     foundAddrs,
			converter: prevFound.converter,
			dataType:  prevFound.dataType,
		})
		if i != len(prevFounds)-1 {
			fmt.Println("------------------------")
		}
	}
	return founds, nil
}

func Attach(pid string) error {
	if isAttached {
		fmt.Println("Already attached.")
		return nil
	}
	fmt.Printf("Target PID: %s\n", pid)

	var err error
	var intPid int
	if intPid, err = strconv.Atoi(pid); err != nil {
		return err
	}

	if err := sys.PtraceAttach(intPid); err == nil {
		fmt.Printf("Attached PID: %s\n", pid)
	} else {
		fmt.Printf("attach failed: %s\n", err)
		return err
	}

	isAttached = true
	return nil
}

func Detach(pid string) error {
	if !isAttached {
		fmt.Println("Already detached.")
		return nil
	}

	var err error
	var intPid int
	if intPid, err = strconv.Atoi(pid); err != nil {
		return err
	}
	if err = sys.PtraceDetach(intPid); err != nil {
		return fmt.Errorf("%s detach failed. %s\n", pid, err)
	} else {
		fmt.Printf("Detached PID: %s\n", pid)
	}

	isAttached = false
	return err
}

func Patch(pid string, targetVal string, targetAddrs []Found) error {
	for _, found := range targetAddrs {
		targetBytes, _ := found.converter(targetVal)
		for _, targetAddr := range found.addrs {
			intPid, _ := strconv.Atoi(pid)
			task := memory.GetTaskForPid(intPid)
			if err := memory.WriteMemory(task, targetAddr, targetBytes); err != nil {
				return err
			}
		}
	}
	fmt.Println("Successfully patched!")
	return nil
}

func Dump(pid string, beginAddress int, endAddress int) error {
	memSize := endAddress - beginAddress
	buf := make([]byte, memSize)
	intPid, _ := strconv.Atoi(pid)
	task := memory.GetTaskForPid(intPid)
	if err := memory.ReadMemory(task, buf, beginAddress, endAddress); err != nil {
		return err
	}
	fmt.Printf("Address range: 0x%x - 0x%x\n", beginAddress, endAddress)
	fmt.Println("--------------------------------------------")
	fmt.Printf("%s", hex.Dump(buf))
	return nil
}
