// Based on Go on iOS(https://github.com/golang/go/blob/master/misc/ios/go_ios_exec.go)
// The license of Go on iOS can be checked in the CREDITS file.

package idevice

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var deviceID string

func Init() (string, error) {
	udid, err := fetchUDID()
	if err != nil {
		return "", err
	}
	deviceID = udid
	return udid, nil
}

func fetchUDID() (string, error) {
	udids := getLines(exec.Command("idevice_id", "-l"))
	if len(udids) == 0 {
		return "", fmt.Errorf("No UDID found; is a device connected?")
	}
	deviceID := string(udids[0])
	return deviceID, nil
}

func getLines(cmd *exec.Cmd) [][]byte {
	out := output(cmd)
	lines := bytes.Split(out, []byte("\n"))
	// Skip the empty line at the end.
	if len(lines[len(lines)-1]) == 0 {
		lines = lines[:len(lines)-1]
	}
	return lines
}

func output(cmd *exec.Cmd) []byte {
	out, err := cmd.Output()
	if err != nil {
		fmt.Println(strings.Join(cmd.Args, "\n"))
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return out
}

func idevCmd(cmd *exec.Cmd) *exec.Cmd {
	if deviceID != "" {
		// Inject -u device_id after the executable, but before the arguments.
		args := []string{cmd.Args[0], "-u", deviceID}
		cmd.Args = append(args, cmd.Args[1:]...)
	}
	return cmd
}

// startDebugBridge ensures that the idevicedebugserverproxy runs on
// port 3222.
func StartDebugBridge() (func(), error) {
	// Kill any hanging debug bridges that might take up port 3222.
	exec.Command("killall", "idevicedebugserverproxy").Run()

	errChan := make(chan error, 1)
	cmd := idevCmd(exec.Command("idevicedebugserverproxy", "-d", "3222"))
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("idevicedebugserverproxy: %v", err)
	}
	go func() {
		if err := cmd.Wait(); err != nil {
			if _, ok := err.(*exec.ExitError); ok {
				errChan <- fmt.Errorf("idevicedebugserverproxy: %s", stderr.Bytes())
			} else {
				errChan <- fmt.Errorf("idevicedebugserverproxy: %v", err)
			}
		}
		errChan <- nil
	}()
	closer := func() {
		cmd.Process.Kill()
		<-errChan
	}
	// Dial localhost:3222 to ensure the proxy is ready.
	delay := time.Second / 4
	for attempt := 0; attempt < 5; attempt++ {
		conn, err := net.DialTimeout("tcp", "localhost:3222", 5*time.Second)
		if err == nil {
			conn.Close()
			return closer, nil
		}
		select {
		case <-time.After(delay):
			delay *= 2
		case err := <-errChan:
			return nil, err
		}
	}
	closer()
	return nil, errors.New("failed to set up idevicedebugserverproxy")
}

// findDevImage use the device iOS version and build to locate a suitable
// developer image.
func findDevImage() (string, error) {
	cmd := idevCmd(exec.Command("ideviceinfo"))
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ideviceinfo: %v", err)
	}
	var iosVer, buildVer string
	lines := bytes.Split(out, []byte("\n"))
	for _, line := range lines {
		spl := bytes.SplitN(line, []byte(": "), 2)
		if len(spl) != 2 {
			continue
		}
		key, val := string(spl[0]), string(spl[1])
		switch key {
		case "ProductVersion":
			iosVer = val
		case "BuildVersion":
			buildVer = val
		}
	}
	if iosVer == "" || buildVer == "" {
		return "", errors.New("failed to parse ideviceinfo output")
	}
	verSplit := strings.Split(iosVer, ".")
	if len(verSplit) > 2 {
		// Developer images are specific to major.minor ios version.
		// Cut off the patch version.
		iosVer = strings.Join(verSplit[:2], ".")
	}
	sdkBase := "/Applications/Xcode.app/Contents/Developer/Platforms/iPhoneOS.platform/DeviceSupport"
	patterns := []string{fmt.Sprintf("%s (%s)", iosVer, buildVer), fmt.Sprintf("%s (*)", iosVer), fmt.Sprintf("%s*", iosVer)}
	for _, pattern := range patterns {
		matches, err := filepath.Glob(filepath.Join(sdkBase, pattern, "DeveloperDiskImage.dmg"))
		if err != nil {
			return "", fmt.Errorf("findDevImage: %v", err)
		}
		if len(matches) > 0 {
			return matches[0], nil
		}
	}
	return "", fmt.Errorf("failed to find matching developer image for iOS version %s build %s", iosVer, buildVer)
}

func MountDevImage() error {
	// Check for existing mount.
	cmd := idevCmd(exec.Command("ideviceimagemounter", "-l", "-x"))
	out, err := cmd.CombinedOutput()
	if err != nil {
		os.Stderr.Write(out)
		return fmt.Errorf("ideviceimagemounter: %v", err)
	}
	var info struct {
		Dict struct {
			Data []byte `xml:",innerxml"`
		} `xml:"dict"`
	}
	if err := xml.Unmarshal(out, &info); err != nil {
		return fmt.Errorf("mountDevImage: failed to decode mount information: %v", err)
	}
	dict, err := parsePlistDict(info.Dict.Data)
	if err != nil {
		return fmt.Errorf("mountDevImage: failed to parse mount information: %v", err)
	}
	if dict["ImagePresent"] == "true" && dict["Status"] == "Complete" {
		return nil
	}
	// Some devices only give us an ImageSignature key.
	if _, exists := dict["ImageSignature"]; exists {
		return nil
	}
	// No image is mounted. Find a suitable image.
	imgPath, err := findDevImage()
	if err != nil {
		return err
	}
	sigPath := imgPath + ".signature"
	cmd = idevCmd(exec.Command("ideviceimagemounter", imgPath, sigPath))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Stderr.Write(out)
		return fmt.Errorf("ideviceimagemounter: %v", err)
	}
	return nil
}

// Parse an xml encoded plist. Plist values are mapped to string.
func parsePlistDict(dict []byte) (map[string]string, error) {
	d := xml.NewDecoder(bytes.NewReader(dict))
	values := make(map[string]string)
	var key string
	var hasKey bool
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if tok, ok := tok.(xml.StartElement); ok {
			if tok.Name.Local == "key" {
				if err := d.DecodeElement(&key, &tok); err != nil {
					return nil, err
				}
				hasKey = true
			} else if hasKey {
				var val string
				var err error
				switch n := tok.Name.Local; n {
				case "true", "false":
					// Bools are represented as <true/> and <false/>.
					val = n
					err = d.Skip()
				default:
					err = d.DecodeElement(&val, &tok)
				}
				if err != nil {
					return nil, err
				}
				values[key] = val
				hasKey = false
			} else {
				if err := d.Skip(); err != nil {
					return nil, err
				}
			}
		}
	}
	return values, nil
}

// findDeviceAppPath returns the device path to the app with the
// given bundle ID. It parses the output of ideviceinstaller -l -o xml,
// looking for the bundle ID and the corresponding Path value.
func FindDeviceAppPath(bundleID string) (string, error) {
	cmd := idevCmd(exec.Command("ideviceinstaller", "-l", "-o", "xml"))
	out, err := cmd.CombinedOutput()
	if err != nil {
		os.Stderr.Write(out)
		return "", fmt.Errorf("ideviceinstaller: -l -o xml %v", err)
	}
	var list struct {
		Apps []struct {
			Data []byte `xml:",innerxml"`
		} `xml:"array>dict"`
	}
	if err := xml.Unmarshal(out, &list); err != nil {
		return "", fmt.Errorf("failed to parse ideviceinstaller output: %v", err)
	}
	for _, app := range list.Apps {
		values, err := parsePlistDict(app.Data)
		if err != nil {
			return "", fmt.Errorf("findDeviceAppPath: failed to parse app dict: %v", err)
		}
		if values["CFBundleIdentifier"] == bundleID {
			if path, ok := values["Path"]; ok {
				return path, nil
			}
		}
	}
	return "", fmt.Errorf("failed to find device path for bundle: %s", bundleID)
}
