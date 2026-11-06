package device

import (
	"time"

	"gitlab.com/CoiaPrant/Sleepy/model"
	pb "gitlab.com/CoiaPrant/Sleepy/proto"
	"gitlab.com/CoiaPrant/cache2go"
	"gitlab.com/CoiaPrant/lib/rand"
)

const defaultDeviceTTL = 20 * time.Second

var devices = cache2go.NewSlidingCache[string, *model.Device]()

func ListDevices() map[string]*model.Device {
	var ret = make(map[string]*model.Device, devices.Size())
	devices.Range(func(session_id string, device *model.Device) (bool, bool) {
		ret[session_id] = device
		return false, false
	})

	return ret
}

func GetDevice(session_id string) (*model.Device, bool) {
	return devices.Touch(session_id)
}

func AddDevice(dev *pb.Device) (session_id string) {
	ttl := max(defaultDeviceTTL, time.Duration(3*dev.ReportInterval))
	device := &model.Device{
		DeviceType:     dev.DeviceType,
		DevicePlatform: dev.DevicePlatform,
		DeviceName:     dev.DeviceName,
		Timestamp:      time.Now().Unix(),
	}

	for {
		session_id = rand.String(20)

		_, loaded := devices.LoadOrStore(session_id, device, ttl)
		if !loaded {
			return
		}
	}
}

func RemoveDevice(session_id string) {
	devices.Delete(session_id)
}
