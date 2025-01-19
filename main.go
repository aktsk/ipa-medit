package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"unsafe"

	"github.com/sterrasec/ipa-medit/pkg/idevice"
	"github.com/sterrasec/ipa-medit/pkg/lldb"
	"github.com/sterrasec/ipa-medit/pkg/prompt"
)

func runApp(binPath string, bundleID string) error {
	// The device app path reported by the device might be stale, so retry
	// the lookup of the device path along with the lldb launching below.
	deviceAppPath, err := idevice.FindDeviceAppPath(bundleID)
	if err != nil {
		return err
	}
	fmt.Printf("Target app: %s\n", deviceAppPath)
	fmt.Printf("Target local bin: %s\n", binPath)
	platform := "remote-ios"
	fmt.Printf("Target platform: %s\n", platform)
	stderr, err := lldb.RunLLDB(platform, binPath, deviceAppPath)
	if err == nil {
		return err
	}
	if len(stderr) != 0 {
		return errors.New(*(*string)(unsafe.Pointer(&stderr)))
	}
	return nil
}

func runMain() error {
	var binPath string
	var bundleID string
	var pid string
	var name string
	flag.StringVar(&binPath, "bin", "", "specify ios app binary that unzip and extract from .ipa")
	flag.StringVar(&bundleID, "id", "", "specify bundle id")
	flag.StringVar(&pid, "pid", "", "specify pid running on the Apple Silicon Mac")
	flag.StringVar(&name, "name", "", "specify process name running on the Apple Silicon Mac")
	flag.Parse()

	if pid != "" {
		fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")
		fmt.Printf("Target PID has been set to %s.\n", pid)
		if result, err := prompt.CheckPidExists(pid); err == nil {
			if result {
				prompt.RunPrompt(pid)
			} else {
				return errors.New("There is no process with the specified pid.")
			}
		} else {
			return err
		}
	}

	if name != "" {
		fmt.Println("Please use `exit` or `Ctrl-D` to exit this program.")
		fmt.Printf("Target Process name has been set to %s.\n", name)
		pid, err := prompt.GetPidByProcessName(name)
		if err != nil {
			return err
		}
		if pid == "" {
			return errors.New("Process not found.")
		}
		prompt.RunPrompt(pid)
		return nil
	}

	if binPath == "" {
		return errors.New("bin option is required.")
	}

	if bundleID == "" {
		return errors.New("id option is required.")
	}

	if err := lldb.PreparePythonFile(); err != nil {
		return err
	}

	udid, err := idevice.Init()
	if err != nil {
		return err
	}
	fmt.Printf("Target device's UDID: %s\n", udid)

	closer, err := idevice.StartDebugBridge()
	if err != nil {
		return err
	}

	fmt.Println("Start to proxy a debugserver connection from a device for remote debugging...")
	defer closer()
	if err := runApp(binPath, bundleID); err != nil {
		return err
	}
	return nil
}

func main() {
	defer prompt.HandleExit()
	err := runMain()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}
