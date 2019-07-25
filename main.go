package main

import (
    "fmt"
    "github.com/JulianSauer/WaterCan/config"
    "github.com/JulianSauer/WaterCan/hue"
    "github.com/JulianSauer/WaterCan/plant_state"
    "github.com/JulianSauer/WaterCan/wireless_sensor_tags"
    "github.com/labstack/echo"
    "net/http"
    "strconv"
    "time"
)

const UPDATE_RATE = 10 * time.Second

func main() {
    fmt.Println("Connecting to Philips Hue...")
    e := hue.InitialConnect()
    if e != nil {
        fmt.Println(e.Error())
    } else {
        fmt.Println("Connected")
    }

    router := echo.New()

    e = wireless_sensor_tags.Login()
    if e != nil {
        router.Logger.Print(e.Error())
    }

    router.GET("forceUpdate", ForceUpdate)
    router.Logger.Fatal(router.Start(":8083"))
}

// Forces an update on all sensors or a specific one
// sensor:  Use an id if only one sensor should be updated
//          If not provided, the default sensors will be loaded from config.json
// timeout: Timeout in seconds
//          If non provided, it will try forever to change the color
func ForceUpdate(context echo.Context) error {
    sensors, e := getSensorFromParam(context)
    if e != nil {
        return context.String(logError(context.Logger(), e))
    }

    timeout, e := getTimeoutFromParam(context)
    if e != nil {
        return context.String(logError(context.Logger(), e))
    }

    lightsToRgb, e := plant_state.ParseValues(sensors)
    if e != nil {
        return context.String(logError(context.Logger(), e))
    }
    for lightId, rgb := range lightsToRgb {
        go updateLight(lightId, rgb[0], rgb[1], rgb[2], timeout, context.Logger())
    }

    return context.String(http.StatusOK, "")
}

func logError(logger echo.Logger, e error) (int, string) {
    logger.Print(e.Error())
    return http.StatusBadRequest, e.Error()
}

func updateLight(id int, red float64, green float64, blue float64, timeout int, logger echo.Logger) {
    timeoutDisabled := false
    if timeout >= 0 {
        logger.Print("Will update light with id ", id, " to color (", red, ", ", green, ", ", blue, ") if it's turned on within the next ", timeout, " seconds")
        timeout = timeout / 10
    } else {
        logger.Print("Will update light with id ", id, " to color (", red, ", ", green, ", ", blue, ") as soon as it's turned on")
        timeoutDisabled = true
    }
    for ; timeoutDisabled || timeout >= 0; timeout-- {
        e := hue.SetLightColor(id, red, green, blue)
        if e != nil {
            if e.Error() == "parameter, xy, is not modifiable. Device is set to off." {
                time.Sleep(UPDATE_RATE)
                continue
            } else {
                logger.Print(e.Error())
            }
        }
        logger.Print("Updated light with id ", id, " to color (", red, ", ", green, ", ", blue, ")")
        break
    }
}

func getSensorFromParam(context echo.Context) ([]*config.Sensor, error) {
    configuration := config.Load()
    if context.QueryParam("sensor") == "" {
        return configuration.Sensors, nil
    } else {
        id, e := strconv.Atoi(context.QueryParam("sensor"))
        if e != nil {
            return nil, e
        }
        return []*config.Sensor{config.GetSensor(configuration.Sensors, id)}, nil
    }
}

func getTimeoutFromParam(context echo.Context) (int, error) {
    if context.QueryParam("timeout") == "" {
        return -1, nil
    } else {
        return strconv.Atoi(context.QueryParam("timeout"))
    }
}
