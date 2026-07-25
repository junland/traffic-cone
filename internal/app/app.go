package app

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"traffic-cone/internal/daemon"
)

const maxHAProxyDataPlanePasswordFileSize = 4 * 1024
const defaultAppName = "traffic-cone"

// Run is the main entry point for the application.
func Run(args []string) int {
	flags := flag.NewFlagSet("traffic-cone", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	pidFile := flags.String("pid-file", filepath.Join(os.TempDir(), fmt.Sprintf("%s.pid", defaultAppName)), "Path to PID file")
	dockerSocket := flags.String("docker-socket", "/var/run/docker.sock", "Path to Docker socket")
	haproxyDataPlaneAPIAddress := flags.String("haproxy-data-plane-api-address", "http://127.0.0.1:5555", "HAProxy Data Plane API service address")
	haproxyDataPlaneAPIUsername := flags.String("haproxy-data-plane-api-username", "", "HAProxy Data Plane API username")
	haproxyDataPlaneAPIPasswordFile := flags.String("haproxy-data-plane-api-password-file", "", "Path to file containing HAProxy Data Plane API password")

	if err := flags.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		return 1
	}

	haproxyDataPlaneAPIPassword, err := resolveHAProxyDataPlanePassword(*haproxyDataPlaneAPIPasswordFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading HAProxy Data Plane API password: %v\n", err)
		return 1
	}

	cfg := daemon.RunConfig{
		AppName:                     defaultAppName,
		PIDFile:                     *pidFile,
		DockerSocket:                *dockerSocket,
		HAProxyDataPlaneAPIAddress:  *haproxyDataPlaneAPIAddress,
		HAProxyDataPlaneAPIUsername: *haproxyDataPlaneAPIUsername,
		HAProxyDataPlaneAPIPassword: haproxyDataPlaneAPIPassword,
	}

	if err := daemon.Start(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting traffic-cone: %v\n", err)
		return 1
	}

	return 0
}

func resolveHAProxyDataPlanePassword(passwordFile string) (string, error) {
	if passwordFile != "" {
		info, err := os.Stat(passwordFile)
		if err != nil {
			return "", err
		}
		if info.Size() > maxHAProxyDataPlanePasswordFileSize {
			return "", fmt.Errorf("password file exceeds %d bytes", maxHAProxyDataPlanePasswordFileSize)
		}

		passwordData, err := os.ReadFile(passwordFile)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(passwordData)), nil
	}
	return os.Getenv("HAPROXY_DATA_PLANE_API_PASSWORD"), nil
}
