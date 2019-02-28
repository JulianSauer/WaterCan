package light

type CapabilitiesControl struct {
    MinDimLevel    int                   `json:"mindimlevel"`
    MaxLumen       int                   `json:"maxlumen"`
    ColorGamutType string                `json:"colorgamuttype"`
    ColorGamut     [][]float32           `json:"colorgamut"`
    Ct             CapabilitiesControlCt `json:"ct"`
}
