package client

import (
	"github.com/godbus/dbus"
	"github.com/muka/device-manager/objects"
	"github.com/muka/device-manager/util"
	"log"
)

// NewDeviceManager create a new DeviceManager client
func NewDeviceManager() *DeviceManager {
	d := new(DeviceManager)
	d.client = NewClient(objects.DeviceManagerInterface, objects.DeviceManagerPath)
	d.logger = util.Logger()
	return d
}

// DeviceManager client
type DeviceManager struct {
	client *Client
	logger *log.Logger
}

// Create a new device
func (d *DeviceManager) Create(dev objects.DeviceDefinition) (dbus.ObjectPath, error) {
	var s dbus.ObjectPath
	err := d.client.Call("Create", 0, dev).Store(&s)
	util.CheckError(err)
	d.logger.Printf("Loaded Device at path\n %v\n", s)
	return s, nil
}

// Read a device definition
func (d *DeviceManager) Read(id string) (*objects.DeviceDefinition, error) {
	var s objects.DeviceDefinition
	err := d.client.Call("Read", 0, id).Store(&s)
	util.CheckError(err)
	return &s, nil
}
