package main

import (
    "github.com/labstack/echo"
    "net/http"
    "strconv"
    "WaterCan/plant_state"
    "WaterCan/hue"
    "fmt"
    "time"
    "WaterCan/wireless_sensor_tags"
    "WaterCan/wireless_sensor_tags/api/logs"
    "math"
    "WaterCan/config"
    "errors"
)

const UPDATE_RATE = 10 * time.Second

func main() {
    e := hue.InitialConnect()
    if e != nil {
        fmt.Println(e.Error())
    }

    router := echo.New()

    e = wireless_sensor_tags.Login()
    if e != nil {
        router.Logger.Print(e.Error())
    }

    router.GET("/update/:sensor", UpdateMoisture)
    router.GET("forceUpdate", ForceUpdate)
    router.Logger.Fatal(router.Start(":8083"))
}

// Updates a light with the given moisture value
// moisture:    Current level that will define the color of a light
// light:       Use the specific id of a Philiphs Hue light
//              If not provided, the default light will be loaded from config.json
func UpdateMoisture(context echo.Context) error {
    moistureLevel, e := getMoistureFromParam(context)
    if e != nil {
        return context.String(logError(context.Logger(), e))
    }

    light, e := getLightFromParam(context)
    if e != nil {
        return context.String(logError(context.Logger(), e))
    }

    rgb := plant_state.Parse(moistureLevel)
    go updateLight(light, rgb[0], rgb[1], rgb[2], context.Logger())

    context.Logger().Print("Current moisture is ", moistureLevel, "%")

    return context.String(http.StatusOK, "")
}

// Forces an update on all sensors or a specific one
// sensor:  Use an id if only one sensor should be updated
//          If not provided, the default sensors will be loaded from config.json
// light:   Use the specific id of a Philiphs Hue light
//          If not provided, the default light will be loaded from config.json
func ForceUpdate(context echo.Context) error {
    ids, e := getSensorFromParam(context)
    if e != nil {
        return context.String(logError(context.Logger(), e))
    }

    moistureLevel, e := moistureOfDriestPlant(ids)
    if e != nil {
        return context.String(logError(context.Logger(), e))
    }

    light, e := getLightFromParam(context)
    if e != nil {
        return context.String(logError(context.Logger(), e))
    }

    rgb := plant_state.Parse(moistureLevel)
    go updateLight(light, rgb[0], rgb[1], rgb[2], context.Logger())

    return context.String(http.StatusOK, "")
}

func logError(logger echo.Logger, e error) (int, string) {
    logger.Print(e.Error())
    return http.StatusBadRequest, e.Error()
}

func updateLight(id int, red float64, green float64, blue float64, logger echo.Logger) {
    logger.Print("Will update light with id ", id, " to color (", red, ", ", green, ", ", blue, ") as soon as it's turned on")
    for {
        e := hue.SetLightColor(id, red, green, blue)
        if e != nil {
            if e.Error() == "parameter, xy, is not modifiable. Device is set to off." {
                time.Sleep(UPDATE_RATE)
                continue
            } else {
                logger.Print(e.Error())
            }
        }
        break
    }
    logger.Print("Updated light with id ", id, " to color (", red, ", ", green, ", ", blue, ")")
}

func moistureOfDriestPlant(ids []int) (float64, error) {
    hourlyStats, e := logs.GetHourlyStats(ids)
    if e != nil {
        return 0, e
    }

    const dateFormat = "1/2/2006"
    var latestDate *time.Time = nil
    var latestStats logs.Stats
    for _, stats := range hourlyStats.Content.Stats {
        date, e := time.Parse(dateFormat, stats.Date)
        if e != nil {
            return 0, e
        }
        if latestDate == nil || latestDate.Before(date) {
            latestDate = &date
            latestStats = stats
        }
    }

    var lowestMoisture = math.MaxFloat64
    for _, values := range latestStats.Values {
        for _, moisture := range values {
            if moisture >= 0 && lowestMoisture > moisture {
                lowestMoisture = moisture
            }
        }
    }
    return lowestMoisture, nil
}

func getMoistureFromParam(context echo.Context) (float64, error) {
    if context.QueryParam("moisture") == "" {
        context.Logger().Print("missing parameter: 'moisture'")
        return 0, errors.New("missing parameter: 'moisture'")
    }

    return strconv.ParseFloat(context.QueryParam("moisture"), 64)
}

func getLightFromParam(context echo.Context) (int, error) {
    if context.QueryParam("light") == "" {
        return config.Load().Light, nil
    } else {
        return strconv.Atoi(context.QueryParam("light"))
    }
}

func getSensorFromParam(context echo.Context) ([]int, error) {
    configuration := config.Load()
    if context.QueryParam("sensor") == "" {
        return configuration.SensorIds, nil
    } else {
        id, e := strconv.Atoi(context.QueryParam("sensor"))
        if e != nil {
            return []int{}, e
        }
        return []int{id}, nil
    }
}
