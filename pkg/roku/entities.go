package roku

type Device struct {
	UDN                string `xml:"udn" json:"udn"`
	SerialNumber       string `xml:"serial-number" json:"serial_number"`
	DeviceID           string `xml:"device-id" json:"device_id"`
	VendorName         string `xml:"vendor-name" json:"vendor_name"`
	ModelName          string `xml:"model-name" json:"model_name"`
	ModelNumber        string `xml:"model-number" json:"model_number"`
	FriendlyDeviceName string `xml:"friendly-device-name" json:"friendly_device_name"`
	Uptime             int64  `xml:"uptime" json:"uptime"`
}

type App struct {
	Name string `xml:"app" json:"app"`
}

type MediaPlayer struct {
	State    string `xml:"state,attr" json:"state"`
	Position string `xml:"position" json:"position"`
}
