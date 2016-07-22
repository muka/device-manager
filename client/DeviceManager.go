package client

import (
	"github.com/godbus/dbus"
	"github.com/muka/device-manager/objects"
)

// NewDeviceManager create a new DeviceManager client
func NewDeviceManager() *DeviceManager {

	d := new(DeviceManager)
	d.client = NewClient(objects.DeviceManagerInterface, objects.DeviceManagerPath)

	return d
}

// DeviceManager client
type DeviceManager struct {
	client *Client
}

// Create a new device
func (d *DeviceManager) Create(dev objects.DeviceDefinition) (dbus.ObjectPath, error) {

	var s dbus.ObjectPath

	err := d.client.Call("Create", 0,
		dev.Id,
		dev.Name,
		dev.Description,
		dev.Path,
		dev.Protocol,
		dev.Properties,
		dev.Streams,
	).Store(&s)
	if err != nil {
		return dbus.ObjectPath("error"), err
	}

	return s, nil
}

// Read a device definition
func (d *DeviceManager) Read(id string) (*objects.DeviceDefinition, error) {

	var s objects.DeviceDefinition
	err := d.client.Call("Read", 0, id).Store(&s)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
