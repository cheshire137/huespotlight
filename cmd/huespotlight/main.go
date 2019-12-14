package main

import (
	"fmt"
	"os"

	"github.com/amimof/huego"
)

const appNameForBridge = "huespotlight"

func getBridgeIPFromUser(bridges []huego.Bridge) (string, error) {
	exitIndex := len(bridges) + 1
	fmt.Println("Found Hue bridges:")
	for i, bridge := range bridges {
		fmt.Printf("%d) %s\n", i+1, bridge.Host)
	}
	fmt.Printf("%d) Cancel\n", exitIndex)
	fmt.Println("Specify which bridge to use:")
	fmt.Print("> ")
	var bridgeIndex int
	_, err := fmt.Scanf("%d", &bridgeIndex)
	if err != nil {
		return "", err
	}
	if bridgeIndex == exitIndex {
		fmt.Println("Exiting...")
		os.Exit(0)
	}
	bridgeIndex--
	if bridgeIndex < 0 || bridgeIndex > len(bridges) {
		return "", fmt.Errorf("Error: %d is not a valid choice, choose between %d and %d",
			bridgeIndex+1, 1, exitIndex)
	}
	return bridges[bridgeIndex].Host, nil
}

func waitForLinkButtonPress(bridge *huego.Bridge) (string, error) {
	buttonPressedChoice := 1
	exitChoice := 2

	fmt.Println("Need to authenticate with Hue bridge, please")
	fmt.Println("press the link button on your bridge...")
	fmt.Printf("%d) Link button has been pressed\n", buttonPressedChoice)
	fmt.Printf("%d) Cancel\n", exitChoice)

	var linkButtonChoice int
	_, err := fmt.Scanf("%d", &linkButtonChoice)
	if err != nil {
		return "", err
	}

	if linkButtonChoice == exitChoice {
		fmt.Println("Exiting...")
		os.Exit(0)
	}
	if linkButtonChoice != buttonPressedChoice {
		return "", fmt.Errorf("Error: %d is not a valid choice, choose between %d and %d",
			linkButtonChoice, buttonPressedChoice, exitChoice)
	}

	newUsername, err := bridge.CreateUser(appNameForBridge)
	if err != nil {
		return "", err
	}

	fmt.Printf("Created user %s\n", newUsername)
	return newUsername, nil
}

func main() {
	bridges, err := huego.DiscoverAll()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var bridgeIP string
	if len(os.Args) < 2 {
		bridgeIP, err = getBridgeIPFromUser(bridges)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		bridgeIP = os.Args[1]
	}

	fmt.Printf("Looking up Hue bridge at IP %s...\n", bridgeIP)
	var bridge *huego.Bridge
	for _, b := range bridges {
		if b.Host == bridgeIP {
			bridge = &b
			break
		}
	}
	if bridge == nil {
		fmt.Println("Error: Could not find bridge")
		os.Exit(1)
	}

	var username string
	if len(os.Args) < 3 {
		username, err = waitForLinkButtonPress(bridge)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		username = os.Args[2]
	}

	fmt.Printf("Logging in as Hue bridge user %s\n", username)
	bridge = bridge.Login(username)

	lights, err := bridge.GetLights()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Found %d lights\n", len(lights))
}
