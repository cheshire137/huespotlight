package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
)

type Config struct {
	AppURLStr           string `json:"app_url"`
	SpotifyAccessToken  string `json:"spotify_access_token"`
	SpotifyRefreshToken string `json:"spotify_refresh_token"`
	SpotifyTokenType    string `json:"spotify_token_type"`
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
	return config, nil
}

func (c *Config) AppURL() (*url.URL, error) {
	serverURL, err := url.Parse(c.AppURLStr)
	if err != nil {
		return nil, err
	}
	return serverURL, nil
}

func (c *Config) SetSpotifyAccessToken(token string) {
	c.SpotifyAccessToken = token
}

func (c *Config) SetSpotifyRefreshToken(token string) {
	c.SpotifyRefreshToken = token
}

func (c *Config) SetSpotifyTokenType(tokenType string) {
	c.SpotifyTokenType = tokenType
}

func (c *Config) Save(path string) error {
	prefix := ""
	indent := "  "
	jsonData, err := json.MarshalIndent(c, prefix, indent)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, jsonData, 0644)
}

func (c *Config) ServerAddr() (string, error) {
	url, err := c.AppURL()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%s", url.Hostname(), url.Port()), nil
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
