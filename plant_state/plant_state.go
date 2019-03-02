package plant_state

import (
    "os"
    "encoding/json"
    "fmt"
)

const CONFIG_NAME = "config.json"

var configFile = load()

type config struct {
    Max float64
}

func Parse(plantState float64) [3]float64 {
    percentage := plantState / configFile.Max
    if percentage >= 0.9 {
        return [3]float64{0, 255, 0}
    } else if percentage >= 0.8 {
        return [3]float64{64, 255, 0}
    } else if percentage >= 0.7 {
        return [3]float64{128, 255, 0}
    } else if percentage >= 0.6 {
        return [3]float64{192, 255, 0}
    } else if percentage >= 0.4 {
        return [3]float64{255, 255, 0}
    } else if percentage >= 0.3 {
        return [3]float64{255, 192, 0}
    } else if percentage >= 0.2 {
        return [3]float64{255, 128, 0}
    } else if percentage >= 0.1 {
        return [3]float64{255, 64, 0}
    } else {
        return [3]float64{255, 0, 0}
    }
}

func load() config {
    file, e := os.Open(CONFIG_NAME)
    if e != nil {
        fmt.Println(e.Error())
    }

    defer file.Close()
    decoder := json.NewDecoder(file)
    config := config{}
    e = decoder.Decode(&config)
    if e != nil {
        fmt.Println("could not parse config.json")
    }
    return config
}
