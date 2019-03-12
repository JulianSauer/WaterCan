package client

import (
    "gopkg.in/resty.v1"
    "encoding/json"
    "WaterCan/wireless_sensor_tags/api"
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
