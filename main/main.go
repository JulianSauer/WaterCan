package main

import (
    "github.com/labstack/echo"
    "net/http"
    "strconv"
)

func main() {
    router := echo.New()
    router.GET("/update/:sensor", UpdateMoisture)
    router.Logger.Fatal(router.Start(":8083"))
}

func UpdateMoisture(context echo.Context) error {
    if context.QueryParam("moisture") == "" {
        context.Logger().Print("Missing parameter: 'moisture'")
        return context.String(http.StatusBadRequest, "Missing parameter: 'moisture'")
    }

    moistureLevel, e := strconv.Atoi(context.QueryParam("moisture"))
    if e != nil {
        context.Logger().Print(e.Error())
        return context.String(http.StatusBadRequest, e.Error())
    }
    sensor := context.Param("sensor")

    context.Logger().Print("Current moisture of "+sensor+" is ", moistureLevel, "%")

    return context.String(http.StatusOK, "")
}
