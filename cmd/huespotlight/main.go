package main

import (
	"fmt"
	"os"

	"github.com/cheshire137/huespotlight/pkg/config"
	"github.com/cheshire137/huespotlight/pkg/hue"
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
}
