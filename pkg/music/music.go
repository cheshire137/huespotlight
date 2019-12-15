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

func (m *Music) GetAuthenticationHandler() func(http.ResponseWriter, *http.Request) {
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

		_, err = w.Write([]byte("Authenticated with Spotify, you can go back to your term now!"))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func (m *Music) GetCurrentSong() error {
	result, err := m.client.PlayerCurrentlyPlaying()
	if err != nil {
		return err
	}
	if !result.Playing {
		return errors.New("nothing currently playing on Spotify")
	}
	artists := make([]string, len(result.Item.Artists))
	for i, artist := range result.Item.Artists {
		artists[i] = artist.Name
	}
	fmt.Printf("Currently playing: %s by %s\n", result.Item.Name, strings.Join(artists, ", "))
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
