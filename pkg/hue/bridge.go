package hue

import (
	"errors"
	"fmt"
	"os"

	"github.com/amimof/huego"

	"github.com/cheshire137/huespotlight/pkg/config"
)

const appNameForBridge = "huespotlight"

// Bridge represents an authenticated connection with a Philips Hue light bright.
type Bridge struct {
	ip     string
	user   string
	client *huego.Bridge
}

// NewBridge constructs a bridge connection by discovering available Philips Hue bridges
// on the network and prompting you to authenticate with one.
func NewBridge(config *config.Config) (*Bridge, error) {
	if len(config.BridgeIP) > 0 && len(config.BridgeUser) > 0 {
		return NewBridgeWithIPAndUser(config.BridgeIP, config.BridgeUser), nil
	}
	if len(config.BridgeIP) > 0 {
		return NewBridgeWithIP(config.BridgeIP)
	}
	bridges, err := huego.DiscoverAll()
	if err != nil {
		return nil, err
	}
	return newBridgeFromList(bridges)
}

// NewBridgeWithIP takes an IP address for a Philips Hue light bridge and returns
// an authenticated connection to that light bridge after prompting the user
// to connect to the bridge.
func NewBridgeWithIP(ip string) (*Bridge, error) {
	bridges, err := huego.DiscoverAll()
	if err != nil {
		return nil, err
	}
	return newBridgeFromListWithIP(bridges, ip)
}

// NewBridgeWithIPAndUser takes an IP address to a Philips Hue bridge and a username
// on that bridge and returns an authenticated connection with that bridge.
func NewBridgeWithIPAndUser(ip string, user string) *Bridge {
	client := huego.New(ip, user)
	client = client.Login(user)
	return &Bridge{ip: ip, user: user, client: client}
}

// FlashLights makes each light on the bridge flash once.
func (b *Bridge) FlashLights() error {
	lights, err := b.client.GetLights()
	if err != nil {
		return err
	}
	for _, light := range lights {
		err = light.Alert("select")
		if err != nil {
			return err
		}
	}
	return nil
}

// TotalLights returns a count of how many lights are registered to the
// bridge.
func (b *Bridge) TotalLights() (int, error) {
	lights, err := b.client.GetLights()
	if err != nil {
		return -1, err
	}
	return len(lights), nil
}

// newBridgeFromListWithIP takes a list of known Philips Hue light bridges and
// an IP address for one of them and returns an authenticated connection to that
// bridge.
func newBridgeFromListWithIP(bridges []huego.Bridge, ip string) (*Bridge, error) {
	fmt.Printf("Looking up Hue bridge at IP %s...\n", ip)
	var client *huego.Bridge
	for _, b := range bridges {
		if b.Host == ip {
			client = &b
			break
		}
	}
	if client == nil {
		return nil, errors.New("Could not find bridge")
	}

	user, err := createUser(client)
	if err != nil {
		return nil, err
	}

	client = client.Login(user)
	return &Bridge{ip: ip, user: user, client: client}, nil
}

// newBridgeFromList takes a list of known Philips Hue light bridges and returns
// an authenticated connection with one of them, based on input from the user.
func newBridgeFromList(bridges []huego.Bridge) (*Bridge, error) {
	ip, err := getIPFromUser(bridges)
	if err != nil {
		return nil, err
	}
	return newBridgeFromListWithIP(bridges, ip)
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

	fmt.Printf("Created bridge user %s\n", newUsername)
	return newUsername, nil
}
