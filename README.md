# Hue Spotlight

A _work in progress_ Go tool to match your Philips Hue lights to your currently playing track on
Spotify.

## How to Run

Create a Spotify application at [developer.spotify.com/my-applications](https://developer.spotify.com/my-applications/).

Edit the app and specify `http://localhost:1234` as your app's redirect URI. This redirect
URI needs to match the address specified in your config.json.

```sh
cp config.json.example config.json
```

Modify config.json to specify the client ID, client secret, and redirect URI for your Spotify app.
If you know the IP address of your Philips Hue bridge, or a user on that bridge, you can specify
those as well. Otherwise, the app will discover bridges on your network and update your config file
to remember the bridge you choose. Then run the app with:

```
go run cmd/huespotlight/main.go config.json
```
