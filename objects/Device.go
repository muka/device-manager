package objects

import (
	"github.com/godbus/dbus"
	"github.com/muka/device-manager/api"
)

//NewDevice create a new Device instance
func NewDevice() (d *Device) {
	return d
}

// Device a device implementation
type Device struct {
	api.Proxy

	Id            string
	Name          string
	Status        int
	Configuration map[string]interface{}
	Profile       map[string]interface{}
	LastUpdate    int
	Data          *RecordObject
	Protocol      *dbus.ObjectPath
}

// Execute an operation
func (d *Device) Execute(op string, payload interface{}) (result ExecuteResult, err *dbus.Error) {
	return result, nil
}

// Connect a device
func (d *Device) Connect() (err *dbus.Error) {
	return nil
}

// Disconnect a device
func (d *Device) Disconnect() (err *dbus.Error) {
	return nil
}

// Read data from a component
func (d *Device) Read(componentId string) (record *RecordObject, err *dbus.Error) {
	return record, nil
}

// Write data to a component
func (d *Device) Write(componentId string, record *RecordObject) (result bool, err *dbus.Error) {
	return result, nil
}

// Subscribe for data updates to a component
func (d *Device) Subscribe(componentId string, params map[string]string) (result bool, err *dbus.Error) {
	return result, nil
}

// Unsubscribe from data updates of a component
func (d *Device) Unsubscribe(componentId string) (result bool, err *dbus.Error) {
	return result, nil
}

// ExecuteResult response object for Execute
type ExecuteResult struct {
	Result     bool
	ResultCode bool
	Response   string
}

// StatusType the status of a device
type StatusType int

const (
	//StatusConnected Device is connected
	StatusConnected StatusType = iota
	//StatusDisconnected Device is not connected
	StatusDisconnected
)
