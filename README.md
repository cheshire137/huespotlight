# Hue Spotlight

A _work in progress_ Go tool to work with your Philips Hue lights.

## How to Run

Create a Spotify application at [developer.spotify.com/my-applications](https://developer.spotify.com/my-applications/). Edit the app and specify a redirect URI.

```sh
cp config.json.example config.json
```

Modify config.json to specify the client ID, client secret, and redirect URI for your Spotify app.
If you know the IP address of your Philips Hue bridge, or a user on that bridge, you can specify
those as well. Otherwise, the app will discover bridges on your network and tell you the IP and user
so you can add them to your config file. Then run the app with:

```
go run cmd/huespotlight/main.go config.json
```
