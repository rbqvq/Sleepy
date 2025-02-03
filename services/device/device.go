package device

import (
	"time"

	"gitlab.com/CoiaPrant/Sleepy/model"
	pb "gitlab.com/CoiaPrant/Sleepy/proto"
	"gitlab.com/CoiaPrant/Sleepy/utils"
	"gitlab.com/CoiaPrant/cache2go"
)

var (
	devices = cache2go.CacheOf[string, *model.Device]()
)

func ListDevices() map[string]*model.Device {
	var devs = make(map[string]*model.Device)
	devices.Foreach(func(session_id string, item *cache2go.CacheItemOf[string, *model.Device]) {
		devs[session_id] = item.Data()
	})
	return devs
}

func GetDevice(session_id string) (*model.Device, error) {
	item, err := devices.Value(session_id)
	if err != nil {
		return nil, err
	}

	return item.Data(), nil
}

func AddDevice(dev *pb.Device) (session_id string) {
	var lifeSpan = 20 * time.Second
	if deviceLifeSpan := time.Duration(dev.ReportInterval * 3); deviceLifeSpan > lifeSpan {
		lifeSpan = deviceLifeSpan
	}

	device := &model.Device{
		DeviceType:     dev.DeviceType,
		DevicePlatform: dev.DevicePlatform,
		DeviceName:     dev.DeviceName,
		Timestamp:      time.Now().Unix(),
	}

	session_id = utils.GetString(20)
	for !devices.NotFoundAdd(session_id, lifeSpan, device) {
		session_id = utils.GetString(20)
	}

	return
}

func RemoveDevice(session_id string) {
	devices.Delete(session_id)
}
