package objects

// BaseQuery base query for devices record
type BaseQuery struct {
	Criteria map[string]string
	OrderBy  map[string]string
	Limit    int32
	Offset   int32
}

// DeviceComponent A device component
type DeviceComponent struct {
	id         string
	unit       string
	format     string
	properties map[string]string
}

// DeviceDefinition A device details list
type DeviceDefinition struct {
	id          string
	name        string
	description string
	path        string
	protocol    string
	properties  map[string]string
	streams     []DeviceComponent
}

// IDeviceManager manages devices in the gateway
type IDeviceManager interface {
	Find(BaseQuery) []IDevice
	Create(dev DeviceDefinition) string
	Read(id string) DeviceDefinition
	Update(id string, dev DeviceDefinition) bool
	Delete(id string) bool
	Batch(operation string, ops map[string]string) bool
}

// DeviceManager manages devices in the gateway
type DeviceManager interface {
	Find(BaseQuery) []IDevice
	Create(dev DeviceDefinition) string
	Read(id string) DeviceDefinition
	Update(id string, dev DeviceDefinition) bool
	Delete(id string) bool
	Batch(operation string, ops map[string]string) bool
}
