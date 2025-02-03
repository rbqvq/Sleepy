package monitor

import (
	"fmt"
	"path"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/shirou/gopsutil/v4/host"
)

var (
	user32   = windows.NewLazyDLL("user32.dll")
	kernel32 = windows.NewLazyDLL("kernel32.dll")

	procGetWindowText             = user32.NewProc("GetWindowTextW")
	procGetWindowTextLength       = user32.NewProc("GetWindowTextLengthW")
	procGetWindowThreadProcessId  = user32.NewProc("GetWindowThreadProcessId")
	procOpenProcess               = kernel32.NewProc("OpenProcess")
	procQueryFullProcessImageName = kernel32.NewProc("QueryFullProcessImageNameW")
)

func getWindowTextLength(hwnd windows.HWND) int {
	ret, _, _ := procGetWindowTextLength.Call(
		uintptr(hwnd))

	return int(ret)
}

func getWindowText(hwnd windows.HWND) string {
	textLen := getWindowTextLength(hwnd) + 1

	buf := make([]uint16, textLen)
	procGetWindowText.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(textLen))

	return syscall.UTF16ToString(buf)
}

func getWindowProcessExecutablePath(hwnd windows.HWND) string {
	pid := getWindowProcessID(hwnd)
	if pid == 0 {
		return "1"
	}

	return getProcessExecutablePath(pid)
}

func getWindowProcessID(hwnd windows.HWND) uint32 {
	var pid uint32
	procGetWindowThreadProcessId.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&pid)),
	)
	return pid
}

func getProcessExecutablePath(pid uint32) string {
	fd, _, _ := procOpenProcess.Call(
		uintptr(windows.PROCESS_QUERY_LIMITED_INFORMATION),
		uintptr(0),
		uintptr(pid),
	)
	if fd == 0 {
		return ""
	}
	defer windows.CloseHandle(windows.Handle(fd))

	buf := make([]uint16, windows.MAX_PATH)
	size := uint32(windows.MAX_PATH)
	procQueryFullProcessImageName.Call(
		fd,
		uintptr(0),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&size)),
	)

	return syscall.UTF16ToString(buf[:size])
}

func GetDeviceType() []string {
	return []string{"windows"}
}

func GetDevicePlatform() string {
	info, err := host.Info()
	if err != nil {
		return "Unknown"
	}

	return fmt.Sprintf("%s %s [%s] (%s)", info.Platform, info.KernelArch, info.PlatformVersion, info.KernelVersion)
}

func GetDeviceName() string {
	info, err := host.Info()
	if err != nil {
		return "Unknown"
	}

	return info.Hostname
}

func GetDeviceState() (using bool, appName string) {
	hwnd := windows.GetForegroundWindow()
	if hwnd == 0 {
		return false, ""
	}

	windowText := getWindowText(hwnd)
	switch windowText {
	case "", "Program Manager":
		return false, ""
	}

	windowExecutable := getWindowProcessExecutablePath(hwnd)
	if windowExecutable == "" {
		return true, windowText
	}

	return true, fmt.Sprintf("[%s] %s", path.Base(strings.ReplaceAll(windowExecutable, "\\", "/")), windowText)
}
