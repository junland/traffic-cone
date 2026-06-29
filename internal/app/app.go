package app

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"traffic-cone/internal/daemon"
)

const appName = "traffic-cone"

func Run(args []string) int {
	if len(args) == 0 {
		printUsage()
		return 1
	}

	if err := runCommand(args[0], args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s failed: %v\n", args[0], err)
		return 1
	}
	return 0
}

func defaultConfig() daemon.RunConfig {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	return daemon.RunConfig{
		AppName:      appName,
		PIDFile:      filepath.Join(cwd, appName+".pid"),
		LogFile:      filepath.Join(cwd, appName+".log"),
		DockerSocket: "/var/run/docker.sock",
		TickInterval: 5 * time.Second,
	}
}

func runCommand(commandName string, args []string) error {
	cfg := defaultConfig()
	cfg.AppName = commandName
	cfg.PIDFile = filepath.Join(filepath.Dir(cfg.PIDFile), commandName+".pid")
	cfg.LogFile = filepath.Join(filepath.Dir(cfg.LogFile), commandName+".log")

	fs := flag.NewFlagSet(commandName, flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	fs.StringVar(&cfg.PIDFile, "pid-file", cfg.PIDFile, "path to pid file")
	fs.StringVar(&cfg.LogFile, "log-file", cfg.LogFile, "path to log file")
	fs.StringVar(&cfg.DockerSocket, "docker-socket", cfg.DockerSocket, "path to docker socket")
	fs.DurationVar(&cfg.TickInterval, "tick", cfg.TickInterval, "heartbeat interval")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if cfg.TickInterval <= 0 {
		return fmt.Errorf("tick must be > 0")
	}
	return daemon.RunForeground(cfg)
}

func printUsage() {
	lines := []string{
		"Generic CLI daemon scaffold",
		"",
		"Usage:",
		"  traffic-cone <daemon-name> [flags]",
		"",
		"Behavior:",
		"  Any daemon name starts that daemon in foreground.",
		"  Use Ctrl+C to stop.",
		"",
		"Examples:",
		"  traffic-cone worker -tick 2s",
		"  traffic-cone scheduler -log-file ./logs/scheduler.log",
	}
	fmt.Println(strings.Join(lines, "\n"))
}
