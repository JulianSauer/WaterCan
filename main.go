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
        router.Logger.Fatal(e.Error())
    }

    router.GET("/update/:sensor/:light", UpdateMoisture)
    router.GET("forceUpdate", ForceUpdate)
    router.Logger.Fatal(router.Start(":8083"))
}

func UpdateMoisture(context echo.Context) error {
    if context.QueryParam("moisture") == "" {
        context.Logger().Print("Missing parameter: 'moisture'")
        return context.String(http.StatusBadRequest, "Missing parameter: 'moisture'")
    }

    moistureLevel, e := strconv.ParseFloat(context.QueryParam("moisture"), 64)
    if e != nil {
        return context.String(logError(context.Logger(), e))
    }

    light, e := strconv.Atoi(context.Param("light"))
    if e != nil {
        return context.String(logError(context.Logger(), e))
    }

    rgb := plant_state.Parse(moistureLevel)
    go updateLight(light, rgb[0], rgb[1], rgb[2], context.Logger())

    context.Logger().Print("Current moisture is ", moistureLevel, "%")

    return context.String(http.StatusOK, "")
}

func ForceUpdate(context echo.Context) error {

    return nil
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
