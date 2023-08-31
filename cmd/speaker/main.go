package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/REPO_DEPRECATED/speaker_alpine/internal/server"
	"github.com/REPO_DEPRECATED/speaker_alpine/internal/translateshell"
)

func main() {
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

	trShell := translateshell.New(ctx)
	go trShell.Run()

	if os.Getenv("WAYLAND_DISPLAY") != "" {
		fmt.Println("Using Wayland")
		err := server.Start(cancel, trShell)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("Using X11")
		// console.Add(cancel, trShell)
		// console.Low()
	}
}
