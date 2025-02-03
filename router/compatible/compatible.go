package compatible

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	. "gitlab.com/CoiaPrant/Sleepy/common/server"
	"gitlab.com/CoiaPrant/Sleepy/model"
	"gitlab.com/CoiaPrant/Sleepy/services/device"
)

var (
	statusOnline = &model.QueryInfo{
		ID:    0,
		Name:  "online",
		Desc:  "Current online",
		Color: "awake",
	}
	statusOffline = &model.QueryInfo{
		ID:    1,
		Name:  "offline",
		Desc:  "Current offline, maybe sleeping",
		Color: "sleeping",
	}
)

func StatusList(c *gin.Context) {
	c.JSON(http.StatusOK, []*model.QueryInfo{statusOnline, statusOffline})
}

func Status(c *gin.Context) {
	now := time.Now().Local()

	info := statusOffline
	deviceList := make(map[string]*model.QueryDevice)
	var lastUpdated int64

	devices := device.ListDevices()
	if len(devices) <= 0 {
		lastUpdated = now.Unix()
	} else {
		for session_id, device := range devices {
			deviceList[session_id] = &model.QueryDevice{
				ShowName: device.DeviceName,
				Using:    device.Using,
				AppName:  device.AppName,
			}

			if device.Using {
				info = statusOnline
			}

			if lastUpdated < device.Timestamp {
				lastUpdated = device.Timestamp
			}
		}
	}

	c.JSON(http.StatusOK, &model.Query{
		Time:              now.Format(time.DateTime),
		Success:           true,
		Status:            info.ID,
		Info:              info,
		Device:            deviceList,
		DeviceStatusSlice: 0,
		LastUpdated:       time.Unix(lastUpdated, 0).Local().Format(time.DateTime),
		Refresh:           REFRESH_INTERVAL.Milliseconds(),
	})
}
