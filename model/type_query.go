package model

type Query struct {
	Time              string                  `json:"time"`
	Success           bool                    `json:"success"`
	Status            int64                   `json:"status"`
	Info              *QueryInfo              `json:"info"`
	Device            map[string]*QueryDevice `json:"device"`
	DeviceStatusSlice int64                   `json:"device_status_slice"`
	LastUpdated       string                  `json:"last_updated"`
	Refresh           int64                   `json:"refresh"`
}

type QueryInfo struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Color string `json:"color"`
}

type QueryDevice struct {
	ShowName string `json:"show_name"`
	Using    bool   `json:"using"`
	AppName  string `json:"app_name"`
}
