package hue

type BridgeSuccess struct {
    Success UserName `json:"success"`
}

type UserName struct {
    Username string `json:"username"`
}
