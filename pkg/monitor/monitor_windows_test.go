package monitor

import (
	"testing"
)

func TestGetDevicePlatform(t *testing.T) {
	t.Logf("device platform: %s", GetDevicePlatform())
}

func TestGetDeviceName(t *testing.T) {
	t.Logf("device name: %s", GetDeviceName())
}

func TestGetDeviceState(t *testing.T) {
	using, appName := GetDeviceState()
	if !using {
		t.Logf("device not using, skip")
		return
	}

	t.Logf("device forceground app name: %s", appName)

}
