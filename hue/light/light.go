package light

type Light struct {
    State            State          `json:"state"`
    SoftwareUpdate   SoftwareUpdate `json:"swupdate"`
    Type             string         `json:"type"`
    Name             string         `json:"name"`
    ModelId          string         `json:"modelid"`
    ManufacturerName string         `json:"manufacturername"`
    ProductName      string         `json:"productname"`
    Capabilities     Capabilities   `json:"capabilities"`
    Config           Config         `json:"config"`
    UniqueId         string         `json:"uniqueid"`
    SoftwareVersion  string         `json:"swversion"`
    SoftwareConfigId string         `json:"swconfigid"`
    ProductId        string         `json:"productid"`
}
