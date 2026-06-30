package app

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"traffic-cone/internal/daemon"
)

// Run is the main entry point for the application.
func Run(args []string) int {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: traffic-cone <daemon-name> [flags]")
		return 1
	}

	daemonName := args[0]
	flags := flag.NewFlagSet(daemonName, flag.ExitOnError)

	pidFile := flags.String("pid-file", filepath.Join(os.TempDir(), fmt.Sprintf("%s.pid", daemonName)), "Path to PID file")
	dockerSocket := flags.String("docker-socket", "unix:///var/run/docker.sock", "Path to Docker socket")
	logFile := flags.String("log-file", filepath.Join(os.TempDir(), fmt.Sprintf("%s.log", daemonName)), "Path to log file")

	if err := flags.Parse(args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		return 1
	}

	cfg := daemon.RunConfig{
		AppName:      daemonName,
		PIDFile:      *pidFile,
		LogFile:      *logFile,
		DockerSocket: *dockerSocket,
	}

	if err := daemon.Start(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting daemon: %v\n", err)
		return 1
	}

	return 0
}
