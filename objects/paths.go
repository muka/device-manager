package objects

const (

	// DeviceManagerPath object path for DeviceManager
	DeviceManagerPath = "/iot/agile/DeviceManager"
	// DeviceManagerInterface interface for DeviceManager
	DeviceManagerInterface = "iot.agile.Device"

	// DevicePath partial path for Device
	DevicePath = "/iot/agile/Device"
	// DeviceInterface interface for Device
	DeviceInterface = "iot.agile.Device"
)

//GetDevicePath return a device path by id
func GetDevicePath(id string) string {
	return DevicePath + "/" + id
}
