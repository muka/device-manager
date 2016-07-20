package objects

import (
	"github.com/godbus/dbus"
	"log"
)

// DeviceManager manages devices in the gateway
type DeviceManager struct {
	BaseObject
	Devices []Device
}

// Find search for devices
func (d *DeviceManager) Find(BaseQuery) (devices []Device, err *dbus.Error) {
	d.Devices = []Device{}
	return d.Devices, err
}

// Create add a device
func (d *DeviceManager) Create(dev DeviceDefinition) (path dbus.ObjectPath, err *dbus.Error) {
	path = "/iot/agile/device/Dummy"
	return path, err
}

// Read read a device definition
func (d *DeviceManager) Read(id string) (dev DeviceDefinition, err *dbus.Error) {

	dev.Id = id
	dev.Description = "My SensorTag device"
	dev.Name = "SensorTag"
	dev.Path = "/iot/agile/device/" + id
	dev.Protocol = "/iot/agile/protocol/BLE"
	dev.Streams = make([]DeviceComponent, 2)

	dev.Streams[0] = DeviceComponent{}
	dev.Streams[0].Id = "temperature"
	dev.Streams[0].Format = "float"
	dev.Streams[0].Unit = "C"

	dev.Streams[1] = DeviceComponent{}
	dev.Streams[1].Id = "light"
	dev.Streams[1].Format = "float"
	dev.Streams[1].Unit = "lumen"

	log.Printf("Read %s: \n%v\n", id, dev)

	return dev, err
}

// Update a device definition
func (d *DeviceManager) Update(id string, dev DeviceDefinition) (res bool, err *dbus.Error) {
	res = true
	return res, err
}

// Delete a device definition
func (d *DeviceManager) Delete(id string) (res bool, err *dbus.Error) {
	res = true
	return res, err
}

// Batch exec batch ops
func (d *DeviceManager) Batch(operation string, ops map[string]string) (res bool, err *dbus.Error) {
	res = true
	return res, err
}

// BaseQuery base query for devices record
type BaseQuery struct {
	Criteria map[string]string
	OrderBy  map[string]string
	Limit    int32
	Offset   int32
}

// DeviceComponent A device component
type DeviceComponent struct {
	Id         string
	Unit       string
	Format     string
	Properties map[string]string
}

// DeviceDefinition A device details list
type DeviceDefinition struct {
	Id          string
	Name        string
	Description string
	Path        string
	Protocol    string
	Properties  map[string]string
	Streams     []DeviceComponent
}
