package config

import (
    "os"
    "fmt"
    "encoding/json"
)

const CONFIG_NAME = "config.json"

type Config struct {
    Max       float64
    Min       float64
    Light     int
    SensorIds []int
}

func Load() *Config {
    file, e := os.Open(CONFIG_NAME)
    if e != nil {
        fmt.Println(e.Error())
    }

    defer file.Close()
    decoder := json.NewDecoder(file)
    config := Config{}
    e = decoder.Decode(&config)
    if e != nil {
        fmt.Println("could not parse config.json")
    }
    return &config
}
