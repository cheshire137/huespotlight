package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
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

	addr, err := config.ServerAddr()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Starting server at %s\n", addr)

	musicClient := music.NewMusic(config)
	handler := musicClient.Authenticate(config)

	if handler != nil {
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

		authDoneChoice := 1
		exitChoice := 2

		fmt.Println("Have you authenticated with Spotify in your browser?")
		fmt.Printf("%d) Yes\n", authDoneChoice)
		fmt.Printf("%d) Cancel\n", exitChoice)

		var userChoice int
		_, err = fmt.Scanf("%d", &userChoice)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		shutdownFunc()

		if userChoice == exitChoice {
			fmt.Println("Exiting...")
			os.Exit(0)
		}
		if userChoice != authDoneChoice {
			fmt.Printf("Error: %d is not a valid choice, choose between %d and %d\n",
				userChoice, authDoneChoice, exitChoice)
			os.Exit(1)
		}

		fmt.Printf("Saving Spotify login information to %s...\n", configPath)
		if err := config.Save(configPath); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	trackID, err := musicClient.GetCurrentSong()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := musicClient.GetSongAnalysis(*trackID); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
