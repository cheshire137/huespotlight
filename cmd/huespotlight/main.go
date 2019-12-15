package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cheshire137/huespotlight/pkg/config"
	"github.com/cheshire137/huespotlight/pkg/hue"
	"github.com/cheshire137/huespotlight/pkg/music"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("%s path_to_config_file.json\n", os.Args[0])
		os.Exit(0)
	}

	configPath := os.Args[1]
	config, err := config.LoadFromFile(configPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Loaded config:")
	fmt.Println(config)

	bridge, err := hue.NewBridge(config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lightCount, err := bridge.TotalLights()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Found %d lights\n", lightCount)
	addr := config.ServerAddr()
	fmt.Printf("Starting server at %s\n", addr)

	musicClient := music.NewMusic(config)
	handler := musicClient.GetAuthenticationHandler()

	server := &http.Server{Addr: addr, Handler: http.HandlerFunc(handler)}
	go func(srv *http.Server) {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Println(err)
			os.Exit(1)
		}
	}(server)

	shutdownFunc := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		fmt.Println("Shutting down server...")
		if err := server.Shutdown(ctx); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	stopSignal := make(chan os.Signal, 1)
	signal.Notify(stopSignal, syscall.SIGINT, syscall.SIGTERM)
	<-stopSignal

	shutdownFunc()
}
