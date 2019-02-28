package hue

type BridgeError struct {
    Error ErrorContent `json:"error"`
}

type ErrorContent struct {
    Type        int    `json:"type"`
    Address     string `json:"address"`
    Description string `json:"description"`
}
