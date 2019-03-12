# WaterCan
As of the current state this simple Go server displays the soil moisture of plants with a Philips Hue light querying [Wireless Sensor Tags](https://store.wirelesstag.net/products/wireless-water-moisture-sensor-2-0) or via a manual update.
Adding some sort of watering mechanism to actually stop plants from dying (and fit the name of the project) is planned but as of right now I lack a water pump.

## Config
The config file defines minimum and maximum values for the moisture levels. Depending on these the color will be more green or red. You can also specify default values for the ids of Wireless Sensor Tags and a Philips Hue light.

Example config:
```
{
  "Max": 25.0,
  "Min": 5.0,
  "Light": 1,
  "SensorIds": [
    1,
    2
  ]
}
```

## REST Endpoints
```
/update/{sensor}?moisture={moisture}&light={light}
```
This sets the color according to `{moisture}` and the minimum/maximum values from config.json. The parameter `{light}` can be set to the id of a light or will default to the value provided in config.json.
If you are using Wireless Sensor Tags you can set an event in their [UI](https://www.mytaglist.com/eth/) under `URL Calling...` and call this from your local network under certain circumstances.

```
/forceUpdated?sensor={sensor}&light={light}
```
This queries all sensors specified in config.json and sets the color of the light matching the state of the driest plant. `{sensor}` can be used to query only one specific sensor and `{light}` can again override the light id from config.json.

# Examples
Suppose the example config from above is used
```
/update/2?moisture=8
```
This will change the color of `light 1` as soon as it's turned on to a red hue since according to the config `5` is the lowest it should go.

```
/update/2?moisture=21&light=2
```
Same as above but this time `light 2` will change to a more green color since the moisture level is still quite high.

```
/forceUpdated
```
Checks all sensors looking for latest data, finds the plant with the lowest moisture level and updates `light 1`.
