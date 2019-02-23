package main

import (
    "github.com/labstack/echo"
    "net/http"
    "strconv"
    "../hue"
    "github.com/amimof/huego"
)

func main() {
    router := echo.New()
    router.GET("/update/:sensor/:light", UpdateMoisture)
    router.Logger.Fatal(router.Start(":8083"))
}

func UpdateMoisture(context echo.Context) error {
    if context.QueryParam("moisture") == "" {
        context.Logger().Print("Missing parameter: 'moisture'")
        return context.String(http.StatusBadRequest, "Missing parameter: 'moisture'")
    }

    moistureLevel, e := strconv.Atoi(context.QueryParam("moisture"))
    if e != nil {
        return context.String(logError(context.Logger(), e))
    }

    light, e := strconv.Atoi(context.Param("light"))
    if e != nil {
        return context.String(logError(context.Logger(), e))
    }

    sensor := context.Param("sensor")

    context.Logger().Print("Current moisture of "+sensor+" is ", moistureLevel, "%")
    bridge, e := hue.Connect()
    bridge.SetLightState(light, huego.State{Hue: 0})

    return context.String(http.StatusOK, "")
}

func logError(logger echo.Logger, e error) (int, string) {
    logger.Print(e.Error())
    return http.StatusBadRequest, e.Error()
}
