package dbus

import (
	"github.com/godbus/dbus/introspect"
	"github.com/muka/device-manager/objects"
)

const (
	objectInterface = "iot.agile.DeviceManager"
	objectPath      = "/iot/agile/DeviceManager"
)

// NewDeviceManagerService create a new instance of DeviceManager DBus API
func NewDeviceManagerService() *DeviceManager {
	dm := &DeviceManager{}
	dm.Object = objects.DeviceManager{}
	dm.Path = objectPath
	dm.Interface = objectInterface
	return dm
}

// DeviceManager daemon instance
type DeviceManager struct {
	Service
	Object objects.DeviceManager
}

// Export the daemon node & intefaces
func (d *DeviceManager) Export() error {

	d.logger.Println("Export to Dbus")

	root := &introspect.Node{
		Name: d.Path,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			{
				Name:    d.Interface,
				Methods: introspect.Methods(&d.Object),
			},
		},
	}

	d.conn.Export(&d.Object,
		d.dbusPath,
		d.Interface)

	d.conn.Export(introspect.NewIntrospectable(root),
		d.dbusPath,
		"org.freedesktop.DBus.Introspectable")

	return nil
}
