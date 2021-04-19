package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"unsafe"

	"github.com/aktsk/ipa-medit/pkg/idevice"
	"github.com/aktsk/ipa-medit/pkg/lldb"
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
	flag.StringVar(&binPath, "bin", "", "specify ios app binary that unzip and extract from .ipa")
	flag.StringVar(&bundleID, "id", "", "specify bundle id")
	flag.Parse()

	if binPath == "" {
		return errors.New("bin option is required")
	}

	if bundleID == "" {
		return errors.New("id option is required")
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

	fmt.Println("Start to proxy a debugserver connection from a device for remote debugging")
	defer closer()
	if err := runApp(binPath, bundleID); err != nil {
		return err
	}
	return nil
}

func main() {
	err := runMain()
	if err != nil {
		log.Fatalf("%v\n", err)
	}
}
