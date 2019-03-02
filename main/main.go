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

    moistureLevel, e := strconv.ParseFloat(context.QueryParam("moisture"), 64)
    if e != nil {
        return context.String(logError(context.Logger(), e))
    }

    light, e := strconv.Atoi(context.Param("light"))
    if e != nil {
        return context.String(logError(context.Logger(), e))
    }

    rgb := plant_state.Parse(moistureLevel)
    e = hue.SetLightColor(2, rgb[0], rgb[1], rgb[2])
    if e != nil {
        return context.String(logError(context.Logger(), e))
    }

    sensor := context.Param("sensor")
    context.Logger().Print("Current moisture of ", sensor, " is ", moistureLevel, "%")
    context.Logger().Print("Updating light with id ", light, " to color (", rgb[0], ", ", rgb[1], ", ", rgb[2], ")")

    return context.String(http.StatusOK, "")
}

func logError(logger echo.Logger, e error) (int, string) {
    logger.Print(e.Error())
    return http.StatusBadRequest, e.Error()
}
