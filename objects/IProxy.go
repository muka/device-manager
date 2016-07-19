package objects

import "github.com/godbus/dbus"

// IProxy generic object
type IProxy interface {
	DBusSetup(conn *dbus.Conn) error
}
