package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	SpotifyClientID string `json:"spotify_client_id"`
	BridgeIP        string `json:"bridge_ip"`
	BridgeUser      string `json:"bridge_user"`
}

func Load() (*Config, error) {
	config := &Config{}
	config.SpotifyClientID = os.Getenv("SPOTIFY_CLIENT_ID")
	config.BridgeIP = os.Getenv("BRIDGE_IP")
	config.BridgeUser = os.Getenv("BRIDGE_USER")
	if err := config.validate(); err != nil {
		return nil, err
	}
	return config, nil
}

func LoadFromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)
	config := &Config{}
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}
	if err := config.validate(); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *Config) String() string {
	return fmt.Sprintf("- Spotify Client ID: %s\n- Philips Hue bridge IP: %s\n- Philips Hue bridge user: %s",
		c.SpotifyClientID, c.BridgeIP, c.BridgeUser)
}

func (c *Config) validate() error {
	if len(c.SpotifyClientID) < 1 {
		return errors.New("no Spotify app client ID in SPOTIFY_CLIENT_ID")
	}
	return nil
}
