package websocket

import (
	"net/http"
	"time"

	. "gitlab.com/CoiaPrant/Sleepy/common/server"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gitlab.com/CoiaPrant/Sleepy/services/device"
)

func Status(c *gin.Context) {
	ws, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.Status(http.StatusUpgradeRequired)
		return
	}
	defer ws.Close()

	var count uint8
	for {
		if err = ws.WriteJSON(device.ListDevices()); err != nil {
			return
		}

		count += 1
		if count%5 == 0 {
			err = ws.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				return
			}
		}

		time.Sleep(time.Second)
	}
}
