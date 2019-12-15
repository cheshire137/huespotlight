package hue

import (
	"errors"
	"fmt"
	"os"

	"github.com/amimof/huego"
)

const appNameForBridge = "huespotlight"

type Bridge struct {
	ip     string
	user   string
	client *huego.Bridge
}

func NewBridge() (*Bridge, error) {
	bridges, err := huego.DiscoverAll()
	if err != nil {
		return nil, err
	}
	return NewBridgeFromList(bridges)
}

func NewBridgeWithIPAndUser(ip string, user string) *Bridge {
	client := huego.New(ip, user)
	return &Bridge{ip: ip, user: user, client: client}
}

func NewBridgeFromList(bridges []huego.Bridge) (*Bridge, error) {
	ip, err := getIPFromUser(bridges)
	if err != nil {
		return nil, err
	}
	return NewBridgeFromListWithIP(bridges, ip)
}

func NewBridgeWithIP(ip string) (*Bridge, error) {
	bridges, err := huego.DiscoverAll()
	if err != nil {
		return nil, err
	}
	return NewBridgeFromListWithIP(bridges, ip)
}

func NewBridgeFromListWithIP(bridges []huego.Bridge, ip string) (*Bridge, error) {
	fmt.Printf("Looking up Hue bridge at IP %s...\n", ip)
	var bridge *huego.Bridge
	for _, b := range bridges {
		if b.Host == ip {
			bridge = &b
			break
		}
	}
	if bridge == nil {
		return nil, errors.New("Could not find bridge")
	}

	user, err := createUser(bridge)
	if err != nil {
		return nil, err
	}

	return &Bridge{ip: ip, user: user, client: bridge}, nil
}

func (b *Bridge) Login() {
	fmt.Printf("Logging in as Hue bridge user %s\n", b.user)
	b.client.Login(b.user)
}

func (b *Bridge) TotalLights() (int, error) {
	lights, err := b.client.GetLights()
	if err != nil {
		return -1, err
	}
	return len(lights), nil
}

func getIPFromUser(bridges []huego.Bridge) (string, error) {
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

func createUser(bridge *huego.Bridge) (string, error) {
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
