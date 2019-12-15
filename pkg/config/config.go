package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
)

type Config struct {
	AppURL              *url.URL
	AppURLStr           string `json:"app_url"`
	SpotifyToken        string `json:"spotify_token"`
	SpotifyClientID     string `json:"spotify_client_id"`
	SpotifyClientSecret string `json:"spotify_client_secret"`
	BridgeIP            string `json:"bridge_ip"`
	BridgeUser          string `json:"bridge_user"`
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
	serverURL, err := url.Parse(config.AppURLStr)
	if err != nil {
		return nil, err
	}
	config.AppURL = serverURL
	return config, nil
}

func (c *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%s", c.AppURL.Hostname(), c.AppURL.Port())
}

func (c *Config) String() string {
	return fmt.Sprintf("- app URL: %s\n- Spotify client ID: %s\n- Philips Hue bridge IP: %s\n- Philips Hue bridge user: %s",
		c.AppURLStr, c.SpotifyClientID, c.BridgeIP, c.BridgeUser)
}

func (c *Config) validate() error {
	if len(c.AppURLStr) < 1 {
		return errors.New("missing app URL")
	}
	if len(c.SpotifyClientID) < 1 {
		return errors.New("missing Spotify app client ID")
	}
	if len(c.SpotifyClientSecret) < 1 {
		return errors.New("missing Spotify app client secret")
	}
	return nil
}
