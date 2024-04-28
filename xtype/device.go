package xtype

//go:generate stringer -type=DevicePlatform -trimprefix=DevicePlatform -output=device.gen.go

type DevicePlatform int16

const (
	DevicePlatformUnknown DevicePlatform = iota
	DevicePlatformIOS
	DevicePlatformAndroid
	DevicePlatformMacOS
	DevicePlatformLinux
	DevicePlatformWindows
	DevicePlatformSafari
	DevicePlatformFirefox
	DevicePlatformChrome

	DevicePlatformCOUNT
)
