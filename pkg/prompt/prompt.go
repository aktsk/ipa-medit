package prompt

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	prompt "github.com/c-bata/go-prompt"

	"github.com/aktsk/ipa-medit/cmd"
)

var appPID string
var addrCache []cmd.Found

func HandleExit() {
	rawModeOff := exec.Command("/bin/stty", "-raw", "echo")
	rawModeOff.Stdin = os.Stdin
	_ = rawModeOff.Run()
	rawModeOff.Wait()
}

func executor(in string) {
	if strings.HasPrefix(in, "find") {
		inputSlice := strings.Split(in, " ")
		dataType := "all"
		targetVal := inputSlice[1]
		if len(inputSlice) < 1 {
			fmt.Println("Target value cannot be specified.")
			return
		}
		if len(inputSlice) == 3 {
			targetVal = inputSlice[2]
			dataType = inputSlice[1]
		}
		if foundAddr, err := cmd.Find(appPID, targetVal, dataType); err == nil {
			addrCache = foundAddr
		}

	} else if strings.HasPrefix(in, "filter") {
		if len(addrCache) == 0 {
			fmt.Println("No previous results. ")
			return
		}
		slice := strings.Split(in, " ")
		if len(slice) == 1 {
			fmt.Println("Target value cannot be specified.")
			return
		}

		foundAddr, err := cmd.Filter(appPID, slice[1], addrCache)
		if err != nil {
			fmt.Println(err)
		}
		addrCache = foundAddr

	} else if strings.HasPrefix(in, "patch") {
		slice := strings.Split(in, " ")
		if len(slice) == 1 {
			fmt.Println("Target value cannot be specified.")
			return
		}

		err := cmd.Patch(appPID, slice[1], addrCache)
		if err != nil {
			fmt.Println(err)
		}

	} else if in == "attach" {
		if err := cmd.Attach(appPID); err != nil {
			HandleExit()
			log.Fatal(err)
		}

	} else if in == "detach" {
		if err := cmd.Detach(appPID); err != nil {
			HandleExit()
			log.Fatal(err)
		}

	} else if strings.HasPrefix(in, "dump") {
		inputSlice := strings.Split(in, " ")
		beginAddr, err := parseAddr(inputSlice[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		endAddr, err := parseAddr(inputSlice[2])
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := cmd.Dump(appPID, beginAddr, endAddr); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

	} else if in == "exit" || in == "quit" {
		fmt.Println("Bye!")
		HandleExit()
		os.Exit(0)

	} else if in == "" {

	} else {
		fmt.Println("Command not found.")
	}
	return
}

func completer(t prompt.Document) []prompt.Suggest {
	return []prompt.Suggest{
		{Text: "find   <int>", Description: "Search the specified integer."},
		{Text: "find   <datatype> <int>", Description: "Types can be specified are string, word, dword, qword."},
		{Text: "filter <int>", Description: "Filter previous search results that match the current search results."},
		{Text: "patch  <int>", Description: "Write the specified value on the address found by search."},
		{Text: "attach", Description: "Attach to the target process by ptrace."},
		{Text: "detach", Description: "Detach from the attached process."},
		{Text: "dump <begin addr> <end addr>", Description: "Display memory dump like hexdump"},
		{Text: "exit"},
	}
}

func RunPrompt(pid string) {
	// for ptrace attach
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	appPID = pid
	p := prompt.New(
		executor,
		completer,
		prompt.OptionTitle("medit: MEmory eDIT tool"),
		prompt.OptionPrefix("> "),
		prompt.OptionInputTextColor(prompt.Cyan),
		prompt.OptionPrefixTextColor(prompt.DarkBlue),
		prompt.OptionPreviewSuggestionTextColor(prompt.Green),
		prompt.OptionDescriptionTextColor(prompt.DarkGray),
	)
	p.Run()
}

func GetPidByProcessName(name string) (string, error) {
	psResult, err := exec.Command("/bin/ps", "-ceo", "pid=,comm=").Output()
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(bytes.NewReader(psResult))
	for scanner.Scan() {
		line := bytes.Split(scanner.Bytes(), []byte(" "))
		lineName := bytes.TrimSpace(line[1])
		linePid := bytes.TrimSpace(line[0])
		if bytes.HasPrefix(lineName, []byte(name)) {
			return string(linePid), nil
		}
	}
	return "", nil
}

func CheckPidExists(pid string) (bool, error) {
	psResult, err := exec.Command("/bin/ps", "-ceo", "pid=,comm=").Output()
	if err != nil {
		return false, err
	}
	scanner := bufio.NewScanner(bytes.NewReader(psResult))
	for scanner.Scan() {
		line := bytes.Split(scanner.Bytes(), []byte(" "))
		lineName := bytes.TrimSpace(line[1])
		linePid := bytes.TrimSpace(line[0])
		if pid == string(linePid) || pid == string(lineName) {
			return true, nil
		}
	}
	return false, nil
}

func parseAddr(arg string) (int, error) {
	arg = strings.Replace(arg, "0x", "", 1)
	address, err := strconv.ParseInt(arg, 16, 64)
	if err == nil {
		return int(address), nil
	}
	address, err = strconv.ParseInt(arg, 10, 64)
	if err == nil {
		return int(address), nil
	}
	return 0, err
}
