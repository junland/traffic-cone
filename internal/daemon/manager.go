package daemon

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// RunConfig controls daemon runtime behavior.
type RunConfig struct {
	AppName      string
	PIDFile      string
	LogFile      string
	DockerSocket string
}

// Start runs the daemon event listening loop in the foreground.
func Start(cfg RunConfig) error {
	if cfg.PIDFile == "" {
		return errors.New("pid file is required")
	}
	if cfg.DockerSocket == "" {
		return errors.New("docker socket is required")
	}
	if cfg.AppName == "" {
		cfg.AppName = "daemon"
	}

	releasePID, err := acquirePIDFile(cfg.PIDFile)
	if err != nil {
		return err
	}
	defer releasePID()

	if err := ensureParentDir(cfg.LogFile); err != nil {
		return fmt.Errorf("prepare log path: %w", err)
	}

	logFile, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "", log.LstdFlags)
	logger.Printf("%s daemon running (pid=%d, docker-socket=%s)", cfg.AppName, os.Getpid(), cfg.DockerSocket)

	ctx := context.Background()
	dockerHost := dockerHostFromSocket(cfg.DockerSocket)
	if !strings.HasPrefix(dockerHost, "unix://") {
		return fmt.Errorf("unsupported docker socket: %s", cfg.DockerSocket)
	}

	socketPath := strings.TrimPrefix(dockerHost, "unix://")
	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
				var dialer net.Dialer
				return dialer.DialContext(ctx, "unix", socketPath)
			},
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://docker/events", nil)
	if err != nil {
		return fmt.Errorf("failed to create Docker events request: %w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect Docker events API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("docker events API returned status %s", resp.Status)
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var msg struct {
			Time   int64  `json:"time"`
			Type   string `json:"Type"`
			Action string `json:"Action"`
			Actor  struct {
				Attributes map[string]string `json:"Attributes"`
			} `json:"Actor"`
		}

		if err := json.Unmarshal(scanner.Bytes(), &msg); err != nil {
			log.Printf("failed to decode docker event: %v\n", err)
			continue
		}

		fmt.Printf("Time: %s | Type: %s | Action: %s | Container: %s | Image: %s\n",
			time.Unix(msg.Time, 0).Format("15:04:05"),
			msg.Type,
			msg.Action,
			msg.Actor.Attributes["name"],
			msg.Actor.Attributes["image"],
		)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("docker events stream read failed: %w", err)
	}

	return nil
}
