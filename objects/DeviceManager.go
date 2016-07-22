package objects

import (
	"log"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/prop"

	"github.com/muka/device-manager/api"
	"github.com/muka/device-manager/service"
	"github.com/muka/device-manager/util"
)

// NewDeviceManager initialize a new DeviceManager object
func NewDeviceManager() *DeviceManager {
	d := DeviceManager{}
	d.Devices = make([]dbus.ObjectPath, 0)
	d.path = DeviceManagerPath
	d.iface = DeviceManagerInterface
	return &d
}

// DeviceManager manages devices in the gateway
type DeviceManager struct {
	api.Proxy

	Devices []dbus.ObjectPath
	devices map[string]*service.Service

	path   string
	iface  string
	logger *log.Logger
}

// GetPath return object path
func (d *DeviceManager) GetPath() string {
	return d.path
}

// SetPath set object path
func (d *DeviceManager) SetPath(s string) {
	d.path = s
}

// GetInterface return interface
func (d *DeviceManager) GetInterface() string {
	return d.iface
}

// SetInterface return interface
func (d *DeviceManager) SetInterface(s string) {
	d.iface = s
}

//SetLogger set default logger
func (d *DeviceManager) SetLogger(logger *log.Logger) {
	d.logger = logger
}

//GetLogger return default logger
func (d *DeviceManager) GetLogger() *log.Logger {
	return d.logger
}

//GetProperties return properties
func (d *DeviceManager) GetProperties() map[string]map[string]*prop.Prop {
	return map[string]map[string]*prop.Prop{
		d.GetInterface(): {
			"Devices": {
				Value:    d.Devices,
				Writable: false,
				Emit:     prop.EmitTrue,
				Callback: func(c *prop.Change) *dbus.Error {
					d.logger.Printf("Changed value %s=%v on %s", c.Name, c.Value, c.Iface)
					return nil
				},
			},
		},
	}
}

// -----
// Dbus API implementation

// Find search for devices
func (d *DeviceManager) Find(q *BaseQuery) (devices []dbus.ObjectPath, err *dbus.Error) {
	d.logger.Println("DeviceManager.Find() not implemented")
	return d.Devices, err
}

// Create add a device
func (d *DeviceManager) Create(dev DeviceDefinition) (path dbus.ObjectPath, err *dbus.Error) {

	id := util.GenerateID()
	d.logger.Printf("Create new device %s\n", id)
	dev.Id = id

	spath := DevicePath + "/" + id
	path = dbus.ObjectPath(spath)
	device := NewDevice(dev)

	service, mErr := service.GetManager().Start(device)
	if mErr != nil {
		d.logger.Fatalf("Cannot start Device service: %s\n", mErr.Error())
		return dbus.ObjectPath("failure"), new(dbus.Error)
	}

	dev.Path = dbus.ObjectPath(device.GetPath())

	d.Devices = append(d.Devices, dev.Path)
	d.devices[dev.Id] = service

	d.logger.Printf("Created Device %s\n", dev.Id)
	return path, err
}

// Read a device definition
func (d *DeviceManager) Read(id string) (dev DeviceDefinition, err *dbus.Error) {

	dev.Id = id
	dev.Description = "My SensorTag device"
	dev.Name = "SensorTag"
	dev.Path = dbus.ObjectPath("/iot/agile/device/Dummy")
	dev.Protocol = dbus.ObjectPath("/iot/agile/protocol/BLE")
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
	Path        dbus.ObjectPath
	Protocol    dbus.ObjectPath
	Properties  map[string]string
	Streams     []DeviceComponent
}
