package roku

type Device struct {
	UDN                string `xml:"udn"`
	SerialNumber       string `xml:"serial-number"`
	DeviceID           string `xml:"device-id"`
	VendorName         string `xml:"vendor-name"`
	ModelName          string `xml:"model-name"`
	ModelNumber        string `xml:"model-number"`
	FriendlyDeviceName string `xml:"friendly-device-name"`
	Uptime             int64  `xml:"uptime"`
}

type ActiveApp struct {
	Name string `xml:"app"`
}

type MediaPlayer struct {
	State    string `xml:"state,attr"`
	Position string `xml:"position"`
}
