package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/back2nix/speaker/internal/config"
	"github.com/back2nix/speaker/internal/localinput"
	"github.com/back2nix/speaker/internal/translateshell"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.Info("Logger initialized")

	cfg, err := config.Init()
	if err != nil {
		logrus.Fatalf("Failed to initialize config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		cancel()
		os.Exit(0)
	}()

	trShell := translateshell.New(ctx, &cfg.Speech)
	go trShell.Run()

	if os.Getenv("WAYLAND_DISPLAY") != "" {
		logrus.Info("Using Wayland input")
	} else if runtime.GOOS == "darwin" { // macOS
		logrus.Info("Using macOS input")
	} else {
		logrus.Info("Using X11 input")
	}
	err = localinput.Start(cancel, trShell, cfg)
	if err != nil {
		logrus.Fatalf("Failed to start local input: %v", err)
	}
}
