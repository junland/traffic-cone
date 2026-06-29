package daemon

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// RunConfig controls daemon runtime behavior.
type RunConfig struct {
	AppName      string
	PIDFile      string
	LogFile      string
	DockerSocket string
	TickInterval time.Duration
}

func Start(cfg RunConfig) error {
	pid, running, err := Status(cfg.PIDFile)
	if err != nil {
		return err
	}
	if running {
		return fmt.Errorf("already running (pid=%d)", pid)
	}

	if err := ensureParentDir(cfg.LogFile); err != nil {
		return fmt.Errorf("prepare log path: %w", err)
	}

	out, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}
	defer out.Close()

	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("resolve executable: %w", err)
	}

	cmd := exec.Command(exe,
		"run",
		"-pid-file", cfg.PIDFile,
		"-log-file", cfg.LogFile,
		"-tick", cfg.TickInterval.String(),
	)
	cmd.Stdout = out
	cmd.Stderr = out

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start daemon process: %w", err)
	}
	_ = cmd.Process.Release()

	deadline := time.Now().Add(4 * time.Second)
	for time.Now().Before(deadline) {
		_, running, _ = Status(cfg.PIDFile)
		if running {
			return nil
		}
		time.Sleep(150 * time.Millisecond)
	}

	return fmt.Errorf("daemon process started but did not become healthy within timeout")
}

func Stop(pidFile string) error {
	pid, running, err := Status(pidFile)
	if err != nil {
		return err
	}
	if !running {
		_ = os.Remove(pidFile)
		return fmt.Errorf("not running")
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("find process: %w", err)
	}
	if err := proc.Kill(); err != nil {
		return fmt.Errorf("stop process %d: %w", pid, err)
	}

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		_, alive, _ := Status(pidFile)
		if !alive {
			_ = os.Remove(pidFile)
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	_ = os.Remove(pidFile)
	return nil
}

func Status(pidFile string) (pid int, running bool, err error) {
	file, err := os.Open(pidFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("read pid file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return 0, false, fmt.Errorf("scan pid file: %w", err)
		}
		return 0, false, fmt.Errorf("pid file is empty")
	}

	pidText := strings.TrimSpace(scanner.Text())
	pid, err = strconv.Atoi(pidText)
	if err != nil {
		return 0, false, fmt.Errorf("invalid pid %q in %s", pidText, pidFile)
	}
	if pid <= 0 {
		return 0, false, fmt.Errorf("invalid pid %d in %s", pid, pidFile)
	}

	return pid, processExists(pid), nil
}

func RunForeground(cfg RunConfig) error {
	releasePID, err := acquirePIDFile(cfg.PIDFile)
	if err != nil {
		return err
	}
	defer releasePID()

	if err := ensureParentDir(cfg.LogFile); err != nil {
		return fmt.Errorf("prepare log path: %w", err)
	}

	logFile, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)
	logger.Printf("%s daemon running (pid=%d, docker-socket=%s)", cfg.AppName, os.Getpid(), cfg.DockerSocket)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ticker := time.NewTicker(cfg.TickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Printf("%s daemon shutdown requested", cfg.AppName)
			return nil
		case t := <-ticker.C:
			logger.Printf("heartbeat at %s", t.Format(time.RFC3339))
		}
	}
}

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

func ensureParentDir(filePath string) error {
	dir := filepath.Dir(filePath)
	if dir == "." || dir == "" {
		return nil
	}
	return os.MkdirAll(dir, 0o755)
}
