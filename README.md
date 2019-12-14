# Hue Spotlight

A _work in progress_ Go tool to work with your Philips Hue lights.

## How to Run

```sh
# Hue bridges on your network will be discovered and you will be able to choose one.
go run cmd/huespotlight/main.go

# Or, if you know the IP address of your Hue bridge, specify it and you will be
# prompted to press the bridge's link button to authenticate:
go run cmd/huespotlight/main.go 192.168.x.x

# Or, if you know the IP address of your Hue bridge and a user on that bridge,
# specify them:
go run cmd/huespotlight/main.go 192.168.x.x ExistingBridgeUserName
```
