package light

type State struct {
    On               bool       `json:"on"`
    Brightness       int        `json:"bri"`
    Hue              int        `json:"hue"`
    Saturation       int        `json:"sat"`
    Effect           string     `json:"effect"`
    Alert            string     `json:"alert"`
    ColorTemperature int        `json:"ct"`
    Xy               [2]float64 `json:"xy"`
    ColorMode        string     `json:"-"`
    Mode             string     `json:"-"`
    Reachable        bool       `json:"-"`
}
