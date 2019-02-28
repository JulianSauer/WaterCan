package hue

import (
    "gopkg.in/resty.v1"
    "errors"
    "strconv"
    "encoding/json"
    "fmt"
    "time"
    "./light"
)

const USER_NAME = "WaterCan"

func InitialConnect() error {
    cache, e := loadCache()
    if e != nil {
        return e
    }

    if cache.User == "" {
        e = register()
        if e != nil {
            return e
        }
    }

    return nil
}

func SetLightState(id int, state *light.State) error {
    baseUrl, e := baseUrl()
    if e != nil {
        return e
    }

    _, e = resty.R().
        SetBody(state).
        Put(baseUrl + "lights/" + strconv.Itoa(id) + "/state")
    return e
}

func GetLightState(id int) (*light.State, error) {
    baseUrl, e := baseUrl()
    if e != nil {
        return nil, e
    }
    response, e := resty.R().Get(baseUrl + "lights/" + strconv.Itoa(id))
    if e != nil {
        return nil, e
    }

    l := light.Light{}
    e = json.Unmarshal(response.Body(), &l)
    if e != nil {
        return nil, e
    }
    return &l.State, nil
}

func register() error {
    cache, e := loadCache()
    if e != nil {
        return e
    }

    var user string
    fmt.Println("Please press the link button on your bridge within the next 20 seconds")
    for i := 0; ; i++ {
        user, e = registerRequest()
        if e == nil || e.Error() == "" {
            break
        }

        if i > 20 {
            return errors.New("link button not pressed")
        }
        time.Sleep(time.Second)
    }

    cache.User = user
    return saveCache(cache)
}

func registerRequest() (string, error) {
    cache, e := loadCache()
    if e != nil {
        return "", e
    }

    response, e := resty.R().
        SetBody([]byte(`{"devicetype":"WaterCan"}`)).
        Post(cache.BridgeIp)
    if e != nil {
        return "", e
    }
    if response.StatusCode() != 200 {
        return "", errors.New("Bad status code: " + strconv.Itoa(response.StatusCode()))
    }
    bridgeError := make([]BridgeError, 0)
    e = json.Unmarshal(response.Body(), &bridgeError)
    if e == nil && bridgeError[0].Error.Description != "" {
        return "", errors.New(bridgeError[0].Error.Description)
    }
    bridgeSuccess := make([]BridgeSuccess, 0)
    e = json.Unmarshal(response.Body(), &bridgeSuccess)
    if e != nil {
        return "", e
    }
    return bridgeSuccess[0].Success.Username, nil
}

func baseUrl() (string, error) {
    cache, e := loadCache()
    if e != nil {
        return "", e
    }
    return cache.BridgeIp + "/" + cache.User + "/", nil
}
