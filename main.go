package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"log/syslog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"gopkg.in/yaml.v3"
)

type Config struct {
	IncludeContainers []string `json:"include_containers" yaml:"include_containers"`
	ExcludeContainers []string `json:"exclude_containers" yaml:"exclude_containers"`
}

var (
	config     Config
	configPath string
	configLock sync.RWMutex
	logger     *log.Logger
)

func main() {
	recreateAll := flag.Bool("a", false, "Recreate all running containers")
	flag.StringVar(&configPath, "c", "", "Path to configuration file")
	flag.Parse()

	// Check for mutually exclusive options
	if *recreateAll && configPath != "" {
		fmt.Println("Error: -a and -c options are mutually exclusive")
		flag.Usage()
		os.Exit(1)
	}

	// Set up syslog logging
	syslogWriter, err := syslog.New(syslog.LOG_INFO|syslog.LOG_DAEMON, "hikup")
	if err != nil {
		log.Fatalf("Error setting up syslog: %v", err)
	}
	logger = log.New(syslogWriter, "", 0)

	// Initial config load if -c is provided
	if configPath != "" {
		if err := reloadConfig(); err != nil {
			logger.Printf("Error loading initial config: %v", err)
			// Continue with default (empty) config
		}
	}

	// Set up signal handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)

	// Start a goroutine to handle SIGHUP
	go func() {
		for {
			<-sigs
			logger.Println("Received SIGHUP, reloading configuration")
			if err := reloadConfig(); err != nil {
				logger.Printf("Error reloading config: %v", err)
			}
		}
	}()

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logger.Fatalf("Error creating Docker client: %v", err)
	}

	for {
		containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
		if err != nil {
			logger.Printf("Error listing containers: %v", err)
			time.Sleep(time.Minute) // Wait before retrying
			continue
		}

		for _, cont := range containers {
			if shouldUpdateContainer(cont, *recreateAll) {
				updateContainer(cli, cont)
			}
		}

		time.Sleep(time.Hour) // Wait for an hour before checking again
	}
}

func reloadConfig() error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading config file: %v", err)
	}

	var newConfig Config
	ext := strings.ToLower(filepath.Ext(configPath))
	switch ext {
	case ".json":
		err = json.Unmarshal(data, &newConfig)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &newConfig)
	default:
		return fmt.Errorf("unsupported config file format: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("error parsing config file: %v", err)
	}

	configLock.Lock()
	config = newConfig
	configLock.Unlock()

	logger.Println("Configuration reloaded successfully")
	return nil
}

func shouldUpdateContainer(cont types.Container, recreateAll bool) bool {
	if recreateAll {
		return true
	}

	configLock.RLock()
	defer configLock.RUnlock()

	// Check if '*' is in the include list
	for _, name := range config.IncludeContainers {
		if name == "*" {
			// Update everything except excluded containers
			return !containsName(config.ExcludeContainers, cont.Names[0][1:])
		}
	}

	// Check if the container is in the include list
	if containsName(config.IncludeContainers, cont.Names[0][1:]) {
		return true
	}

	// Check if the container is in the exclude list
	if containsName(config.ExcludeContainers, cont.Names[0][1:]) {
		return false
	}

	// If not in include or exclude list, don't update by default
	return false
}

func containsName(names []string, target string) bool {
	for _, name := range names {
		if name == target {
			return true
		}
	}
	return false
}

func updateContainer(cli *client.Client, cont types.Container) {
	// ... (rest of the updateContainer function remains the same)
}
