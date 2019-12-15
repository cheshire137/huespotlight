package main

import (
	"fmt"
	"os"

	"github.com/cheshire137/huespotlight/pkg/hue"
)

func main() {
	var bridge *hue.Bridge
	var err error

	if len(os.Args) > 2 {
		bridgeIP := os.Args[1]
		username := os.Args[2]
		bridge = hue.NewBridgeWithIPAndUser(bridgeIP, username)
	} else if len(os.Args) > 1 {
		bridgeIP := os.Args[1]
		bridge, err = hue.NewBridgeWithIP(bridgeIP)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		bridge, err = hue.NewBridge()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	lightCount, err := bridge.TotalLights()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Found %d lights\n", lightCount)
}
