// Based on Go on iOS(https://github.com/golang/go/blob/master/misc/ios/go_ios_exec.go)
// The license of Go on iOS can be checked in the CREDITS file.

package lldb

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"reflect"
	"syscall"
)

var pyPath string

func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}

func calcFileHash(filepath string) []byte {
	f, _ := os.Open(filepath)
	defer f.Close()

	h := sha256.New()
	io.Copy(h, f)

	return h.Sum(nil)
}

func makePythonFile(filepath string) error {
	fp, err := os.Create(filepath)
	if err != nil {
		return err
	}
	fp.WriteString(pythonData)
	return nil
}

func PreparePythonFile() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	pyPath = filepath.Join(filepath.Dir(exePath), "lldb-driver.py")
	if !fileExists(pyPath) {
		makePythonFile(pyPath)
	} else {
		pyDataHash := sha256.Sum256([]byte(pythonData))
		fileHash := calcFileHash(pyPath)
		if !reflect.DeepEqual(pyDataHash[:], fileHash) {
			makePythonFile(pyPath)
		}
	}
	return nil
}

func RunLLDB(platform, binPath, deviceAppPath string) ([]byte, error) {
	var env []string
	for _, e := range os.Environ() {
		env = append(env, e)
	}
	lldb := exec.Command(
		"xcrun",
		"python3",
		pyPath,
		platform,
		binPath,
		deviceAppPath,
	)
	lldb.Env = env
	lldb.Stdin = os.Stdin
	lldb.Stdout = os.Stdout
	var stderr bytes.Buffer
	lldb.Stderr = io.MultiWriter(&stderr, os.Stderr)
	err := lldb.Start()
	if err == nil {
		// Forward SIGQUIT to the lldb driver which in turn will forward
		// to the running program.
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGQUIT)
		proc := lldb.Process
		go func() {
			for sig := range sigs {
				proc.Signal(sig)
			}
		}()
		err = lldb.Wait()
		signal.Stop(sigs)
		close(sigs)
	}
	return stderr.Bytes(), err
}
