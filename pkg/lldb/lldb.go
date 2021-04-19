package lldb

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
)

var pyPath string

func fileExists(filepath string) bool {
    _, err := os.Stat(filepath)
    return os.IsExist(err)
}

func PreparePythonFile() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	pyPath = filepath.Join(filepath.Dir(exePath), "lldb-driver.py")
	if !fileExists(pyPath) {
		fmt.Println("aaa")
		fp, err := os.Create(pyPath)
		if err != nil {
			return err
		}
		fp.WriteString(pythonData)
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
