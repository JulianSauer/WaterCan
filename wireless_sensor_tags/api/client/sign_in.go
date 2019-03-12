package client

import (
    "gopkg.in/resty.v1"
    "WaterCan/wireless_sensor_tags/api"
)

type SignInRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func SignIn(username string, password string) error {
    body := SignInRequest{username, password}

    response, e := resty.R().
        SetBody(body).
        Post(api.CLIENT_URL + "SignIn")
    if e != nil {
        return e
    }

    return api.ParseError(response)
}
