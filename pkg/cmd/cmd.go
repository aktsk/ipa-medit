package cmd

import (
	//"errors"
	"fmt"
	"os/exec"

	"github.com/aktsk/ipa-medit/pkg/memory"
)

type Found struct {
	addrs     []int
	converter func(string) ([]byte, error)
	dataType  string
}

func Find(pid string) ([]Found, error) {
	founds := []Found{}
	out, err := exec.Command("vmmap", "--wide", pid).Output()
	if err != nil {
		return nil, err
	}
	addrRanges, _ := memory.GetWritableAddrRanges(out)
	fmt.Println(addrRanges)
	return founds, nil

	//return founds, errors.New("Error: specified datatype does not exist")
}