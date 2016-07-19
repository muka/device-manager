package objects

import (
	"fmt"
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/prop"
)

// DeviceManager manages devices in the gateway
type DeviceManager struct {
	IDeviceManager
	Devices []Device
}

// Find search for devices
func (d *DeviceManager) Find(BaseQuery) []Device {
	d.Devices = []Device{}
	return d.Devices
}

// Create add a device
func (d *DeviceManager) Create(dev DeviceDefinition) string {
	return "ciao"
}

// Read read a device definition
func (d *DeviceManager) Read(id string) DeviceDefinition {
	return DeviceDefinition{}
}

// Update a device definition
func (d *DeviceManager) Update(id string, dev DeviceDefinition) bool {
	return true
}

// Delete a device definition
func (d *DeviceManager) Delete(id string) bool {
	return true
}

// Batch exec batch ops
func (d *DeviceManager) Batch(operation string, ops map[string]string) bool {
	return true
}

// DBusSetup the DeviceManager integration to dbus
func (d *DeviceManager) DBusSetup(conn *dbus.Conn) error {

	propsSpec := map[string]map[string]*prop.Prop{
		"iot.agile.DeviceManager": {
			"Devices": {
				[]string{},
				true,
				prop.EmitTrue,
				func(c *prop.Change) *dbus.Error {
					fmt.Println(c.Name, "changed to", c.Value)
					return nil
				},
			},
		},
	}

	// props :=
	prop.New(conn, "/iot/agile/DeviceManager", propsSpec)

	return nil
}
