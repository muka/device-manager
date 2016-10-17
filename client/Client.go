package client

import (
	"log"

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

	c.Config = config
	c.logger = logger

	return c
}

// Client implement a DBus client
type Client struct {
	logger     *log.Logger
	conn       *dbus.Conn
	dbusObject dbus.BusObject
	Config     *Config
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
		if c.Config.Bus == SystemBus {
			c.logger.Println("Using SystemBus")
			return dbus.SystemBus()
		} else if c.Config.Bus == SessionBus {
			c.logger.Println("Using SessionBus")
			return dbus.SessionBus()
		} else {
			c.logger.Println("Unknown Bus!")
			return nil, nil
		}
	}

	dbusConn, err := getConn()
	if err != nil {
		return err
	}

	c.conn = dbusConn
	c.dbusObject = c.conn.Object(c.Config.Name, dbus.ObjectPath(c.Config.Path))

	c.logger.Printf("Connected to %s %s\n", c.Config.Name, c.Config.Path)
	return nil
}

// Call a DBus method
func (c *Client) Call(method string, flags dbus.Flags, args ...interface{}) *dbus.Call {

	if !c.isConnected() {
		err := c.Connect()
		if err != nil {
			return &dbus.Call{
				Err: err,
			}
		}
	}

	methodPath := c.iface + "." + method

	callArgs := args
	c.logger.Printf("Call %s( %v )\n", methodPath, callArgs)

	return c.dbusObject.Call(methodPath, flags, callArgs...)
}

//GetProperty return a property value
func (c *Client) GetProperty(p string) (dbus.Variant, error) {
	if !c.isConnected() {
		return dbus.Variant{}, c.Connect()
	}
	return c.dbusObject.GetProperty(p)
}

//GetProperties load all the properties for an interface
func (c *Client) GetProperties(props interface{}) error {

	if !c.isConnected() {
		err := c.Connect()
		if err != nil {
			return err
		}
	}

	c.logger.Printf("Loading properties for %s", c.Config.Iface)

	result := make(map[string]dbus.Variant)
	err := c.dbusObject.Call("org.freedesktop.DBus.Properties.GetAll", 0, c.Config.Iface).Store(&result)
	if err != nil {
		return err
	}

	return util.MapToStruct(props, result)
}
