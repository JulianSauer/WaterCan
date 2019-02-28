package light

type Capabilities struct {
    Certified bool                  `json:"certified"`
    Control   CapabilitiesControl   `json:"control"`
    Streaming CapabilitiesStreaming `json:"streaming"`
}
