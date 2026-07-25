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
	flags := flag.NewFlagSet("traffic-cone", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	pidFile := flags.String("pid-file", filepath.Join(os.TempDir(), fmt.Sprintf("%s.pid", daemonName)), "Path to PID file")
	dockerSocket := flags.String("docker-socket", "/var/run/docker.sock", "Path to Docker socket")
	haproxyDataPlaneAPIAddress := flags.String("haproxy-data-plane-api-address", "http://127.0.0.1:5555", "HAProxy Data Plane API service address")
	haproxyDataPlaneAPIUsername := flags.String("haproxy-data-plane-api-username", "", "HAProxy Data Plane API username")
	haproxyDataPlaneAPIPassword := flags.String("haproxy-data-plane-api-password", "", "HAProxy Data Plane API password")

	if err := flags.Parse(args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		return 1
	}

	cfg := daemon.RunConfig{
		AppName:                        daemonName,
		PIDFile:                        *pidFile,
		DockerSocket:                   *dockerSocket,
		HAProxyDataPlaneAPIAddress:     *haproxyDataPlaneAPIAddress,
		HAProxyDataPlaneAPIUsername:    *haproxyDataPlaneAPIUsername,
		HAProxyDataPlaneAPIPassword:    *haproxyDataPlaneAPIPassword,
	}

	if err := daemon.Start(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting traffic-cone: %v\n", err)
		return 1
	}

	return 0
}
