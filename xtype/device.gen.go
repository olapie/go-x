// Code generated by "stringer -type=DevicePlatform -trimprefix=DevicePlatform -output=device.gen.go"; DO NOT EDIT.

package xtype

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[DevicePlatformUnknown-0]
	_ = x[DevicePlatformIOS-1]
	_ = x[DevicePlatformAndroid-2]
	_ = x[DevicePlatformMacOS-3]
	_ = x[DevicePlatformLinux-4]
	_ = x[DevicePlatformWindows-5]
	_ = x[DevicePlatformSafari-6]
	_ = x[DevicePlatformFirefox-7]
	_ = x[DevicePlatformChrome-8]
	_ = x[DevicePlatformCOUNT-9]
}

const _DevicePlatform_name = "UnknownIOSAndroidMacOSLinuxWindowsSafariFirefoxChromeCOUNT"

var _DevicePlatform_index = [...]uint8{0, 7, 10, 17, 22, 27, 34, 40, 47, 53, 58}

func (i DevicePlatform) String() string {
	if i < 0 || i >= DevicePlatform(len(_DevicePlatform_index)-1) {
		return "DevicePlatform(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _DevicePlatform_name[_DevicePlatform_index[i]:_DevicePlatform_index[i+1]]
}