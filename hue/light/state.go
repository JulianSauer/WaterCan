package light

type State struct {
    On               bool   `json:"on"`
    Brightness       int    `json:"bri"`
    Hue              int    `json:"hue"`
    Saturation       int    `json:"sat"`
    Effect           string `json:"effect"`
    Alert            string `json:"alert"`
    ColorTemperature int    `json:"ct"`
    // Not modifiable:
    Xy        []float32 `json:"-"`
    ColorMode string    `json:"-"`
    Mode      string    `json:"-"`
    Reachable bool      `json:"-"`
}
