package config

import (
    "encoding/json"
    "fmt"
    "os"
)

const CONFIG_NAME = "config.json"

type Config struct {
    Sensors []*Sensor
}

type Sensor struct {
    Max      float64
    Min      float64
    Light    int
    SensorId int
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

func GetSensor(sensors []*Sensor, sensorId int) *Sensor {
    for _, sensor := range sensors {
        if sensor.SensorId == sensorId {
            return sensor
        }
    }
    return nil
}

func GetSensors(sensors []*Sensor, sensorIds []int) []*Sensor {
    var truncatedSensors []*Sensor
    for _, sensor := range sensors {
        for _, sensorId := range sensorIds {
            if sensor.SensorId == sensorId {
                truncatedSensors = append(truncatedSensors, sensor)
                break
            }
        }
    }
    return truncatedSensors
}
