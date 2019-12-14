package main

import (
	"fmt"
	"os"

	"github.com/amimof/huego"
	"github.com/cheshire137/huespotlight/pkg/hue"
)

func main() {
	var bridge *hue.Bridge

	if len(os.Args) > 2 {
		bridgeIP := os.Args[1]
		username := os.Args[2]
		bridge = hue.NewBridge(bridgeIP, username)
	} else {
		bridges, err := huego.DiscoverAll()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if len(os.Args) > 1 {
			bridgeIP := os.Args[1]
			bridge, err = hue.NewBridgeFromListWithIP(bridges, bridgeIP)
		} else {
			bridge, err = hue.NewBridgeFromList(bridges)
		}
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	bridge.Login()

	lightCount, err := bridge.TotalLights()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Found %d lights\n", lightCount)
}
