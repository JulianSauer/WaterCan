package hue

import (
    "github.com/amimof/huego"
    "os"
    "encoding/json"
    "fmt"
)

type cache struct {
    BridgeIp string
}

const USER = "WaterCan"
const CACHE = "cache.json"

var bridge *huego.Bridge

func Connect() (*huego.Bridge, error) {
    if bridge != nil {
        return bridge, nil
    }

    bridgeIp, e := loadBridgeIP()
    if e != nil || bridgeIp == "" {
        e = initialConnect()
        if e != nil {
            return nil, nil
        }
    } else {
        bridge = huego.New(bridgeIp, USER)
    }

    fmt.Print("Connected to " + bridgeIp)
    return bridge, nil
}

func initialConnect() error {
    var e error
    bridge, e = huego.Discover()
    if e != nil {
        return e
    }

    user, e := bridge.CreateUser(USER)
    if e != nil {
        return e
    }

    bridge = bridge.Login(user)
    fmt.Println("Please press the button on your Philips Hue Bridge")
    return saveBridgeIP()
}

func loadBridgeIP() (string, error) {
    file, e := os.Open(CACHE)
    if e != nil {
        return "", e
    }

    defer file.Close()
    decoder := json.NewDecoder(file)
    cache := cache{}
    e = decoder.Decode(&cache)
    return cache.BridgeIp, e
}

func saveBridgeIP() error {
    cache := cache{BridgeIp: bridge.Host}
    file, e := os.Create(CACHE)
    if e != nil {
        return e
    }

    encoder := json.NewEncoder(file)
    return encoder.Encode(cache)
}
