package daemon

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/moby/moby/client"
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
	if cfg.LogFile == "" {
		return errors.New("log file is required")
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

	// Initialize Docker client
	cli, err := client.New(client.WithHost(cfg.DockerSocket))
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()

	ctx := context.Background()

	// Subscribe to the docker events API stream
	result := cli.Events(ctx, client.EventsListOptions{})
	messages := result.Messages
	errs := result.Err

	// Block and listen to incoming events asynchronously
	for {
		select {
		case err := <-errs:
			if err != nil {
				log.Printf("Event stream error: %v\n", err)
			}
		case msg := <-messages:
			// Handle the specific event trigger here
			fmt.Printf("Time: %s | Type: %s | Action: %s | Container: %s | Image: %s\n",
				time.Unix(msg.Time, 0).Format("15:04:05"),
				msg.Type,
				msg.Action,
				msg.Actor.Attributes["name"],
				msg.Actor.Attributes["image"],
			)
		}
	}
}
