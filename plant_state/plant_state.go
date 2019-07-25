package plant_state

import (
    "github.com/JulianSauer/WaterCan/config"
    "github.com/JulianSauer/WaterCan/wireless_sensor_tags/api/logs"
    "time"
    "math"
)

func ParseValues(sensors []*config.Sensor) (map[int][3]float64, error) {
    lightToRgb := make(map[int][3]float64) // light -> rgb

    // Collect all sensors that use the same light
    sensorsOfLights := make(map[int][]*config.Sensor) // light -> sensors
    for _, sensor := range sensors {
        sensorsOfLights[sensor.Light] = append(sensorsOfLights[sensor.Light], sensor)
    }

    for lightId, sensors := range sensorsOfLights {
        rgb, e := moistureAsRgbOfDriestPlant(sensors)
        if e != nil {
            return nil, e
        }
        lightToRgb[lightId] = rgb
    }
    return lightToRgb, nil
}

// Compares the hourly stats of every given sensor using a percentage value
// See plant_state#moistureToPercentage
func moistureAsRgbOfDriestPlant(sensors []*config.Sensor) ([3]float64, error) {
    sensorsToHourlyStats := make(map[*config.Sensor]*logs.HourlyStats)
    for _, sensor := range sensors {
        if hourlyStats, e := logs.GetHourlyStats([]int{sensor.SensorId}); e == nil {
            sensorsToHourlyStats[sensor] = hourlyStats
        } else {
            return [3]float64{0, 255, 0}, e
        }
    }

    const dateFormat = "1/2/2006"
    var latestStats logs.Stats
    var lowestPercentage = math.MaxFloat64
    for sensor, hourlyStats := range sensorsToHourlyStats {
        var latestDate *time.Time = nil
        for _, stats := range hourlyStats.Content.Stats {
            date, e := time.Parse(dateFormat, stats.Date)
            if e != nil {
                return [3]float64{0, 255, 0}, e
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

        percentage := moistureToPercentage(sensor, lowestMoisture)
        if percentage < lowestPercentage {
            lowestPercentage = percentage
        }
    }

    return percentageToRgb(lowestPercentage), nil
}

func moistureToPercentage(sensor *config.Sensor, plantState float64) float64 {
    return (plantState - sensor.Min) /
        (sensor.Max - sensor.Min)
}

func percentageToRgb(percentage float64) [3]float64 {
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
