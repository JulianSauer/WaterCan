package hue

import (
    "gopkg.in/resty.v1"
    "errors"
    "strconv"
    "encoding/json"
    "fmt"
    "time"
    "./light"
    "math"
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

func SetLightColor(id int, red float64, green float64, blue float64) error {
    xy := ConvertRGBToXY(red, green, blue)
    body := fmt.Sprintf("{\"xy\": [%f, %f]}", xy[0], xy[1])
    return SetLightState(id, body)
}

func SetLightState(id int, body string) error {
    baseUrl, e := baseUrl()
    if e != nil {
        return e
    }

    response, e := resty.R().
        SetBody(body).
        Put(baseUrl + "lights/" + strconv.Itoa(id) + "/state")
    bridgeError := make([]BridgeError, 0)
    e = json.Unmarshal(response.Body(), &bridgeError)
    if e == nil && bridgeError[0].Error.Description != "" {
        return errors.New(bridgeError[0].Error.Description)
    }
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

// See https://stackoverflow.com/a/22649803
func ConvertRGBToXY(red float64, green float64, blue float64) [2]float64 {
    var normalizedToOne [3]float64
    normalizedToOne[0] = red / 255
    normalizedToOne[1] = green / 255
    normalizedToOne[2] = blue / 255
    red = enhanceColor(normalizedToOne[0])
    green = enhanceColor(normalizedToOne[1])
    blue = enhanceColor(normalizedToOne[2])

    x := red*0.649926 + green*0.103455 + blue*0.197109
    y := red*0.234327 + green*0.743075 + blue*0.022598
    z := red*0.0000000 + green*0.053077 + blue*1.035763

    return [2]float64{
        x / (x + y + z),
        y / (x + y + z)}
}

// See https://stackoverflow.com/a/22649803
func enhanceColor(normalizedColor float64) float64 {
    if normalizedColor > 0.04045 {
        return math.Pow((normalizedColor+0.055)/(1.0+0.055), 2.4)
    } else {
        return normalizedColor / 12.92
    }
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
