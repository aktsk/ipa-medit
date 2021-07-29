package prompt

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
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
		{Text: "ps", Description: "Find the target process and if there is only one, specify it as the target."},
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
