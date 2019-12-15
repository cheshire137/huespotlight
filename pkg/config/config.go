package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	SpotifyClientID    string `json:"spotify_client_id"`
	SpotifyRedirectURI string `json:"spotify_redirect_uri"`
	BridgeIP           string `json:"bridge_ip"`
	BridgeUser         string `json:"bridge_user"`
}

func Load() (*Config, error) {
	config := &Config{}
	config.SpotifyClientID = os.Getenv("SPOTIFY_CLIENT_ID")
	config.SpotifyRedirectURI = os.Getenv("SPOTIFY_REDIRECT_URI")
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
	return fmt.Sprintf("- Spotify client ID: %s\n- Spotify redirect URI: %s\n- Philips Hue bridge IP: %s\n- Philips Hue bridge user: %s",
		c.SpotifyClientID, c.SpotifyRedirectURI, c.BridgeIP, c.BridgeUser)
}

func (c *Config) validate() error {
	if len(c.SpotifyClientID) < 1 {
		return errors.New("missing Spotify app client ID")
	}
	if len(c.SpotifyRedirectURI) < 1 {
		return errors.New("missing Spotify app redirect URI")
	}
	return nil
}
