package plant_state

import (
    "os"
    "encoding/json"
    "fmt"
)

const CONFIG_NAME = "config.json"
const GREEN float32 = 25500

var configFile = load()

type config struct {
    Max float32
}

func Parse(plantState float32) int {
    return int((plantState / configFile.Max) * GREEN)
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
