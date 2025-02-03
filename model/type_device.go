package model

type Device struct {
	DeviceType     []string `json:"device_type"`
	DevicePlatform string   `json:"device_platform"`
	DeviceName     string   `json:"device_name"`
	AppName        string   `json:"app_name"`
	Using          bool     `json:"using"`
	Timestamp      int64    `json:"timestamp"`
}
