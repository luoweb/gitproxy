package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/git-proxy/go/internal/config"
	"github.com/git-proxy/go/internal/logger"
	"github.com/git-proxy/go/internal/proxy"
	"github.com/git-proxy/go/internal/ssh"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	if *showVersion {
		fmt.Println("git-proxy-go v1.0.0")
		return
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	if *configPath != "" {
		dir := path.Dir(*configPath)
		if dir != "." {
			os.Chdir(dir)
		}
	}

	log := logger.New(cfg.Logging.Level, cfg.Logging.Format)

	gitProxy := proxy.New(cfg, log)

	if err := gitProxy.Start(); err != nil {
		log.Error("Failed to start git proxy: %v", err)
		os.Exit(1)
	}

	sshProxy := ssh.New(cfg, log)
	if err := sshProxy.Start(); err != nil {
		log.Error("Failed to start SSH proxy: %v", err)
		gitProxy.Stop()
		os.Exit(1)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Info("Received signal: %s", sig)

	sshProxy.Stop()
	gitProxy.Stop()
}
