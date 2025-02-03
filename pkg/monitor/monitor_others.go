//go:build !windows

package monitor

import "runtime"

func GetDeviceType() []string {
	return []string{runtime.GOOS}
}

func GetDevicePlatform() string {
	return "Unknown"
}

func GetDeviceName() string {
	return "Unknown"
}

func GetDeviceState() (using bool, appName string) {
	return true, "Unknown"
}
