package client

import (
	"github.com/muka/device-manager/objects"
	"github.com/muka/device-manager/util"
	"log"
)

// NewDevice create a new Device client
func NewDevice(id string) *Device {
	d := new(Device)
	d.client = NewClient(objects.DeviceInterface, objects.DevicePath+"/"+id)
	d.logger = util.Logger()
	return d
}

// Device client
type Device struct {
	client *Client
	logger *log.Logger
}

// Read data from a device definition
func (d *Device) Read(componentId string) (*objects.RecordObject, error) {
	var s objects.RecordObject
	err := d.client.Call("Read", 0, componentId).Store(&s)
	return &s, err
}
