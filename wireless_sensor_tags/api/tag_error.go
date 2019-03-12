package api

import (
    "encoding/json"
    "gopkg.in/resty.v1"
    "errors"
)

type TagError struct {
    Message       string `json:"Message"`
    Stacktrace    string `json:"Stacktrace"`
    ExceptionType string `json:"ExceptionType"`
}

func ParseError(response *resty.Response) error {
    var tagError = TagError{}
    e := json.Unmarshal(response.Body(), &tagError)
    if e == nil && tagError.Message != "" {
        return errors.New(tagError.Message)
    }
    return e
}
