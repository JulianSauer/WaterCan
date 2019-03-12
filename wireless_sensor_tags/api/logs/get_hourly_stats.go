package logs

import (
    "gopkg.in/resty.v1"
    "WaterCan/wireless_sensor_tags/api"
    "encoding/json"
)

type GetHourlyStatsRequest struct {
    Ids  []int  `json:"ids"`
    Type string `json:"type"`
}

type HourlyStats struct {
    Content Content `json:"d"`
}

type Content struct {
    Type     string   `json:"__type"`
    Stats    []Stats  `json:"stats"`
    TempUnit int      `json:"temp_unit"`
    Ids      []int    `json:"ids"`
    Names    []string `json:"names"`
}

type Stats struct {
    Date   string      `json:"date"`
    Ids    []int       `json:"ids"`
    Values [][]float64 `json:"values"`
}

func GetHourlyStats(ids []int) (*HourlyStats, error) {
    body := GetHourlyStatsRequest{ids, "cap"}

    response, e := resty.R().
        SetBody(body).
        Post(api.LOGS_URL + "GetHourlyStats")

    if e != nil {
        return nil, e
    }
    e = api.ParseError(response)
    if e != nil {
        return nil, e
    }

    hourlyStats := HourlyStats{}
    e = json.Unmarshal(response.Body(), &hourlyStats)

    return &hourlyStats, e
}
