package client

import (
	"log"
	// "reflect"

	// "github.com/fatih/structs"
	"github.com/godbus/dbus"
	"github.com/muka/device-manager/util"
)

// NewClient create a new client
func NewClient(iface string, path string) *Client {

	c := new(Client)

	logger, err := util.NewLogger("client")
	util.CheckError(err)

	c.logger = logger
	c.path = path
	c.iface = iface

	return c
}

// Client implement a DBus client
type Client struct {
	logger     *log.Logger
	conn       *dbus.Conn
	dbusObject dbus.BusObject
	iface      string
	path       string
}

func (c *Client) isConnected() bool {
	return c.conn != nil
}

// Connect connect to DBus
func (c *Client) Connect() error {

	c.logger.Println("Connecting to DBus")

	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}

	c.conn = conn
	c.dbusObject = conn.Object(c.iface, dbus.ObjectPath(c.path))

	c.logger.Println("Connected")

	return nil
}

// Call a DBus method
func (c *Client) Call(method string, flags dbus.Flags, args ...interface{}) *dbus.Call {

	if !c.isConnected() {
		c.Connect()
	}

	methodPath := c.iface + "." + method

	// var callArgs = make([]interface{}, 0)
	// for _, arg := range args {
	// 	if reflect.TypeOf(arg).Kind() == reflect.Struct {
	//
	// 		for _, value := range structs.Values(arg) {
	// 			callArgs = append(callArgs, value)
	// 		}
	//
	// 	} else {
	// 		callArgs = append(callArgs, arg)
	// 	}
	// }

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
