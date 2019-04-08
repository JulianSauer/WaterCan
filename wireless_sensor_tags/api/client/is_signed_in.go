package client

import (
    "encoding/json"
    "github.com/JulianSauer/WaterCan/wireless_sensor_tags/api"
    "gopkg.in/resty.v1"
)

type IsSignedInResponse struct {
    IsSignedIn bool `json:"d"`
}

func IsSignedIn() (bool, error) {
    response, e := resty.R().
        Post(api.CLIENT_URL + "IsSignedIn")
    if e != nil {
        return false, e
    }

    var isSignedInResponse IsSignedInResponse
    e = json.Unmarshal(response.Body(), &isSignedInResponse)
    return isSignedInResponse.IsSignedIn, nil
}
