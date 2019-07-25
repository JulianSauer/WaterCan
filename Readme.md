# WaterCan
As of the current state this simple Go server displays the soil moisture of plants with a Philips Hue light querying [Wireless Sensor Tags](https://store.wirelesstag.net/products/wireless-water-moisture-sensor-2-0).
Adding some sort of watering mechanism to actually stop plants from dying (and fit the name of the project) is planned but as of right now I lack a water pump.

## Config
The config file defines minimum and maximum values for the moisture levels.
Depending on these the color of the specified light will be more green or red.
If you use the id of a light for more than one sensor the driest plant in relation to it's minimum/maximum values will be picked for display. 

Example config:
```
{
  "Sensors": [
    {
      "Max": 25.0,
      "Min": 5.0,
      "Light": 1,
      "SensorId": 1
    },
    {
      "Max": 30.0,
      "Min": 8.0,
      "Light": 1,
      "SensorId": 2
    },
    {
      "Max": 20.0,
      "Min": 3.0,
      "Light": 2,
      "SensorId": 3
    }
  ]
}
```

## REST Endpoints
```
/forceUpdated?sensor={sensorId}&timeout={timeoutInSeconds}
```
This queries all sensors specified in config.json and sets the color of their lights matching the state of the driest plant.
`{sensorId}` can be used to query only one specific sensor and `{timeoutInSeconds}` defines how long the program tries to update lights if they are currently turned off.
If no timeout is provided it will try to update until the lamp is turned on eventually.

# Examples
Suppose the example config from above is used
```
/forceUpdated
```
Checks all sensors looking for latest data.
`sensor 3` will use `light 2` to display the current state.
Since `sensor 1` and `sensor 2` both use `light 1`, their moisture will be calculated in relation to their minimum and maximum values and `light 1` will display the lower value.  

```
/forceUpdated?sensor=2&timeout=3600
```
The moisture level of `sensor 2` will be displayed using `light 1`.
If it is currently turned of the updated will be retried for a maximum time of one hour.
