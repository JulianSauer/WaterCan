package main

import (
    "github.com/labstack/echo"
    "net/http"
    "strconv"
    "../plant_state"
    "../hue"
    "fmt"
)

func main() {
    e := hue.InitialConnect()
    if e != nil {
        fmt.Println(e.Error())
    }

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

    hueState := plant_state.Parse(float32(moistureLevel))
    state, e := hue.GetLightState(2)
    if e != nil {
        return context.String(logError(context.Logger(), e))
    }
    state.Hue = hueState
    e = hue.SetLightState(2, state)
    if e != nil {
        return context.String(logError(context.Logger(), e))
    }

    sensor := context.Param("sensor")
    context.Logger().Print("Current moisture of ", sensor, " is ", moistureLevel, "%")
    context.Logger().Print("Updating light with id ", light, " to color ", hueState)

    return context.String(http.StatusOK, "")
}

func logError(logger echo.Logger, e error) (int, string) {
    logger.Print(e.Error())
    return http.StatusBadRequest, e.Error()
}
