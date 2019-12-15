package music

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"

	"github.com/cheshire137/huespotlight/pkg/config"
)

type Music struct {
	auth   *spotify.Authenticator
	client *spotify.Client
}

func NewMusic(config *config.Config) *Music {
	auth := spotify.NewAuthenticator(config.AppURLStr, spotify.ScopeUserReadPlaybackState)
	auth.SetAuthInfo(config.SpotifyClientID, config.SpotifyClientSecret)
	return &Music{auth: &auth}
}

func (m *Music) Authenticate(config *config.Config) func(http.ResponseWriter, *http.Request) {
	if len(config.SpotifyAccessToken) > 0 {
		fmt.Println("Using Spotify token from config to authenticate...")
		token := &oauth2.Token{AccessToken: config.SpotifyAccessToken}
		if len(config.SpotifyRefreshToken) > 0 {
			token.RefreshToken = config.SpotifyRefreshToken
		}
		if len(config.SpotifyTokenType) > 0 {
			token.TokenType = config.SpotifyTokenType
		}
		client := m.auth.NewClient(token)
		m.client = &client
		return nil
	}

	state := getRandomString(10)
	authURL := m.auth.AuthURL(state)
	fmt.Println("Please visit this URL to authenticate:")
	fmt.Printf("\t%s\n", authURL)

	return func(w http.ResponseWriter, r *http.Request) {
		token, err := m.auth.Token(state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusNotFound)
			return
		}

		fmt.Println("Authenticated with Spotify")
		fmt.Printf("- Spotify token: %s\n", token.AccessToken)
		client := m.auth.NewClient(token)
		m.client = &client

		config.SetSpotifyAccessToken(token.AccessToken)
		config.SetSpotifyRefreshToken(token.RefreshToken)
		config.SetSpotifyTokenType(token.TokenType)

		_, err = w.Write([]byte("Authenticated with Spotify, you can go back to your term now!"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func (m *Music) GetCurrentSong() (*spotify.ID, error) {
	result, err := m.client.PlayerCurrentlyPlaying()
	if err != nil {
		return nil, err
	}
	if !result.Playing {
		return nil, errors.New("nothing currently playing on Spotify")
	}
	artists := make([]string, len(result.Item.Artists))
	for i, artist := range result.Item.Artists {
		artists[i] = artist.Name
	}
	fmt.Printf("Currently playing: %s by %s\n", result.Item.Name, strings.Join(artists, ", "))
	return &result.Item.ID, nil
}

func (m *Music) GetSongAnalysis(id spotify.ID) error {
	analysis, err := m.client.GetAudioAnalysis(id)
	if err != nil {
		return err
	}
	fmt.Println("Beats")
	fmt.Println(analysis.Beats)
	return nil
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func getRandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
