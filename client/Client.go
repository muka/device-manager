package client

import (
	"log"
	// "reflect"

	// "github.com/fatih/structs"
	"github.com/godbus/dbus"
	"github.com/muka/device-manager/util"
)

const (
	// SessionBus uses the session bus
	SessionBus = 0
	// SystemBus uses the system bus
	SystemBus = 1
)

// Config pass configuration to a DBUS client
type Config struct {
	Name  string
	Iface string
	Path  string
	Bus   int
}

// NewClient create a new client
func NewClient(config *Config) *Client {

	c := new(Client)

	logger, err := util.NewLogger("client")
	util.CheckError(err)

	c.logger = logger

	c.path = config.Path
	c.iface = config.Iface
	c.name = config.Name
	c.bus = config.Bus

	return c
}

// Client implement a DBus client
type Client struct {
	logger     *log.Logger
	conn       *dbus.Conn
	dbusObject dbus.BusObject
	bus        int
	iface      string
	path       string
	name       string
}

func (c *Client) isConnected() bool {
	return c.conn != nil
}

// Connect connect to DBus
func (c *Client) Connect() error {

	c.logger.Println("Connecting to DBus")

	var getConn = func() (*dbus.Conn, error) {
		if c.bus == SystemBus {
			c.logger.Println("Using SystemBus")
			return dbus.SystemBus()
		} else if c.bus == SessionBus {
			c.logger.Println("Using SessionBus")
			return dbus.SessionBus()
		} else {
			return nil, nil
		}
	}

	dbusConn, err := getConn()
	if err != nil {
		return err
	}

	c.conn = dbusConn
	c.dbusObject = c.conn.Object(c.name, dbus.ObjectPath(c.path))

	c.logger.Printf("Connected to %s %s\n", c.name, c.path)

	return nil
}

// Call a DBus method
func (c *Client) Call(method string, flags dbus.Flags, args ...interface{}) *dbus.Call {

	if !c.isConnected() {
		c.Connect()
	}

	methodPath := c.iface + "." + method

	callArgs := args
	c.logger.Printf("Call %s( %v )\n", methodPath, callArgs)

	return c.dbusObject.Call(methodPath, flags, callArgs...)
}

//GetProperty return a property value
func (c *Client) GetProperty(p string) (dbus.Variant, error) {
	if !c.isConnected() {
		c.Connect()
	}
	return c.dbusObject.GetProperty(p)
}
