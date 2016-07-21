package dbus

import (
	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"
	"github.com/muka/device-manager/objects"
)

const (
	objectInterface = "iot.agile.DeviceManager"
	objectPath      = "/iot/agile/DeviceManager"
)

// NewDeviceManagerService create a new instance of DeviceManager DBus API
func NewDeviceManagerService() *DeviceManager {
	dm := &DeviceManager{}
	dm.Object = objects.NewDeviceManager()
	dm.Path = objectPath
	dm.Interface = objectInterface
	return dm
}

// DeviceManager daemon instance
type DeviceManager struct {
	Service
	Object objects.DeviceManager
}

func (d *DeviceManager) exportProperties() (*prop.Properties, error) {

	propsSpec := map[string]map[string]*prop.Prop{
		d.Interface: {
			"Devices": {
				d.Object.Devices,
				false, // Writable
				prop.EmitTrue,
				func(c *prop.Change) *dbus.Error {
					d.logger.Println(c.Name, "changed to", c.Value)
					return nil
				},
			},
		},
	}

	props := prop.New(d.conn, d.dbusPath, propsSpec)
	return props, nil
}

// Export the daemon node & intefaces
func (d *DeviceManager) Export() error {

	d.logger.Println("Export to Dbus")

	// Expose the objerct path
	d.conn.Export(&d.Object,
		d.dbusPath,
		d.Interface)

	// Register properties
	props, err := d.exportProperties()
	if err != nil {
		return err
	}

	root := &introspect.Node{
		Name: d.Path,
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			{
				Name:       d.Interface,
				Methods:    introspect.Methods(&d.Object),
				Properties: props.Introspection(d.Interface),
			},
		},
	}

	// Export Introspectable
	d.conn.Export(introspect.NewIntrospectable(root),
		d.dbusPath,
		"org.freedesktop.DBus.Introspectable")

	// c := make(chan *dbus.Signal)
	// d.conn.Signal(c)
	// for _ = range c {
	// }
	return nil
}
