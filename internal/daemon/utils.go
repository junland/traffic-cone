package daemon

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// ensureParentDir ensures that the parent directory of the given file path exists.
func ensureParentDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}

// Status reads the PID file and checks if the process is running.
func Status(pidFile string) (pid int, running bool, err error) {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("read pid file: %w", err)
	}

	pidText := strings.TrimSpace(string(data))
	if pidText == "" {
		return 0, false, fmt.Errorf("pid file is empty")
	}
	if nl := strings.IndexByte(pidText, '\n'); nl >= 0 {
		pidText = strings.TrimSpace(pidText[:nl])
	}

	pid, err = strconv.Atoi(pidText)
	if err != nil {
		return 0, false, fmt.Errorf("invalid pid %q in %s", pidText, pidFile)
	}
	if pid <= 0 {
		return 0, false, fmt.Errorf("invalid pid %d in %s", pid, pidFile)
	}

	return pid, processExists(pid), nil
}

// aqcuirePIDFile creates a PID file and returns a function to release it.
func acquirePIDFile(pidFile string) (func(), error) {
	pid, running, err := Status(pidFile)
	if err != nil {
		var pathErr *os.PathError
		if !errors.As(err, &pathErr) {
			return nil, err
		}
	}
	if running {
		return nil, fmt.Errorf("already running (pid=%d)", pid)
	}

	if err := ensureParentDir(pidFile); err != nil {
		return nil, fmt.Errorf("prepare pid path: %w", err)
	}

	pidData := []byte(strconv.Itoa(os.Getpid()) + "\n")
	if err := os.WriteFile(pidFile, pidData, 0o644); err != nil {
		return nil, fmt.Errorf("write pid file: %w", err)
	}

	return func() {
		_ = os.Remove(pidFile)
	}, nil
}

// dockerHostFromSocket converts a socket path to a Docker host URL.
func dockerHostFromSocket(socketPath string) string {
	if strings.Contains(socketPath, "://") {
		return socketPath
	}
	if strings.HasPrefix(socketPath, "//./pipe/") {
		return "npipe://" + socketPath
	}
	return "unix://" + socketPath
}

// processExists checks if a process with the given PID exists.
func processExists(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))
	return err == nil
}
