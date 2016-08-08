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

// Find a list of device
func (d *DeviceManager) Find(q *objects.BaseQuery) ([]dbus.ObjectPath, error) {
	var s []dbus.ObjectPath
	err := d.client.Call("Find", 0, q).Store(&s)
	return s, err
}

// Create a new device
func (d *DeviceManager) Create(dev objects.DeviceDefinition) (objects.DeviceDefinition, error) {
	var s objects.DeviceDefinition
	err := d.client.Call("Create", 0, dev).Store(&s)
	return s, err
}

// Update a device
func (d *DeviceManager) Update(id string, dev objects.DeviceDefinition) (bool, error) {
	var s bool
	err := d.client.Call("Create", 0, id, dev).Store(&s)
	return s, err
}

// Read a device definition
func (d *DeviceManager) Read(id string) (*objects.DeviceDefinition, error) {
	var s objects.DeviceDefinition
	err := d.client.Call("Read", 0, id).Store(&s)
	return &s, err
}

// Delete a device definition
func (d *DeviceManager) Delete(id string) (bool, error) {
	var s bool
	err := d.client.Call("Delete", 0, id).Store(&s)
	return s, err
}
