package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	ENUM_CURRENT_SETTINGS  = ^uint(0)
	CDS_UPDATEREGISTRY     = 0x01
	CDS_TEST               = 0x02
	DISP_CHANGE_SUCCESSFUL = 0
	DISP_CHANGE_RESTART    = 1
	DISP_CHANGE_FAILED     = ^uintptr(0)
	DMDO_DEFAULT           = 0
	DMDO_90                = 1
	DMDO_180               = 2
	DMDO_270               = 3
)

type DEVMODE struct {
	dmDeviceName         [32]uint16
	dmSpecVersion        int16
	dmDriverVersion      int16
	dmSize               int16
	dmDriverExtra        int16
	dmFields             int32
	dmPositionX          int32
	dmPositionY          int32
	dmDisplayOrientation int32
	dmDisplayFixedOutput int32
	dmColor              int16
	dmDuplex             int16
	dmYResolution        int16
	dmTTOption           int16
	dmCollate            int16
	dmFormName           [32]uint16
	dmLogPixels          int16
	dmBitsPerPel         int16
	dmPelsWidth          int32
	dmPelsHeight         int32
	dmDisplayFlags       int32
	dmDisplayFrequency   int32
	dmICMMethod          int32
	dmICMIntent          int32
	dmMediaType          int32
	dmDitherType         int32
	dmReserved1          int32
	dmReserved2          int32
	dmPanningWidth       int32
	dmPanningHeight      int32
}

var (
	user32                = syscall.NewLazyDLL("user32.dll")
	enumDisplaySettings   = user32.NewProc("EnumDisplaySettingsW")
	changeDisplaySettings = user32.NewProc("ChangeDisplaySettingsW")
)

func main() {
	var dm DEVMODE
	dm.dmSize = int16(unsafe.Sizeof(dm))
	if ret, _, _ := enumDisplaySettings.Call(0, uintptr(ENUM_CURRENT_SETTINGS), uintptr(unsafe.Pointer(&dm))); ret != 0 {
		temp := dm.dmPelsHeight
		dm.dmPelsHeight = dm.dmPelsWidth
		dm.dmPelsWidth = temp

		switch dm.dmDisplayOrientation {
		case DMDO_DEFAULT:
			dm.dmDisplayOrientation = DMDO_90
		case DMDO_270:
			dm.dmDisplayOrientation = DMDO_180
		case DMDO_180:
			dm.dmDisplayOrientation = DMDO_90
		case DMDO_90:
			dm.dmDisplayOrientation = DMDO_DEFAULT
		}

		if ret, _, _ := changeDisplaySettings.Call(uintptr(unsafe.Pointer(&dm)), uintptr(CDS_TEST)); syscall.Handle(ret) == syscall.Handle(DISP_CHANGE_FAILED) {
			fmt.Println("resolution not supported")
		} else {
			if ret, _, _ := changeDisplaySettings.Call(uintptr(unsafe.Pointer(&dm)), uintptr(CDS_UPDATEREGISTRY)); syscall.Handle(ret) == syscall.Handle(DISP_CHANGE_SUCCESSFUL) {
				fmt.Println("resolution changed successfully")
			} else if syscall.Handle(ret) == syscall.Handle(DISP_CHANGE_RESTART) {
				fmt.Println("you need to reboot for the change to happen")
			} else {
				fmt.Println("failed to change resolution")
			}
		}
	} else {
		fmt.Println("failed to get current resolution")
	}
}
