package objects

import (
	"log"

	"github.com/godbus/dbus/prop"

	"github.com/godbus/dbus"
	"github.com/muka/device-manager/api"
	"github.com/robertkrimen/otto"
)

//NewDevice create a new Device instance
func NewDevice(def DeviceDefinition) *Device {

	dev := Device{}

	dev.Id = def.Id
	dev.Name = def.Name
	dev.Status = StatusDisconnected
	dev.Configuration = make(map[string]string)

	spath := DevicePath + "/" + dev.Id
	dev.SetPath(spath)
	dev.SetInterface(DeviceInterface)

	dev.Profile = make(map[string]DeviceComponent)
	for _, stream := range def.Streams {
		dev.Profile[stream.Id] = stream
	}

	dev.LastUpdate = 0
	dev.Data = &RecordObject{}
	dev.Protocol = &def.Protocol

	return &dev
}

// Device a device implementation
type Device struct {
	api.Proxy

	logger *log.Logger
	path   string
	iface  string

	Id            string
	Name          string
	Status        StatusType
	Configuration map[string]string
	Profile       map[string]DeviceComponent
	LastUpdate    int32
	Data          *RecordObject
	Protocol      *dbus.ObjectPath
}

// GetPath return object path
func (d *Device) GetPath() string {
	return d.path
}

// SetPath set object path
func (d *Device) SetPath(s string) {
	d.path = s
}

// GetInterface return interface
func (d *Device) GetInterface() string {
	return d.iface
}

// SetInterface return interface
func (d *Device) SetInterface(s string) {
	d.iface = s
}

//SetLogger set default logger
func (d *Device) SetLogger(logger *log.Logger) {
	d.logger = logger
}

//GetLogger return default logger
func (d *Device) GetLogger() *log.Logger {
	return d.logger
}

//GetProperties return properties
func (d *Device) GetProperties() map[string]map[string]*prop.Prop {
	return map[string]map[string]*prop.Prop{
		d.GetInterface(): {
			"Data": {
				Value:    d.Data,
				Writable: false,
				Emit:     prop.EmitTrue,
				Callback: func(c *prop.Change) *dbus.Error {
					d.logger.Printf("Device.Data: Changed value %s=%v on %s", c.Name, c.Value, c.Iface)
					return nil
				},
			},
		},
	}
}

// -----
// Dbus API implementation

// Execute an operation
func (d *Device) Execute(op string, payload string) (result ExecuteResult, err *dbus.Error) {
	d.logger.Print("Device.Execute() not implemented")
	return result, err
}

// Connect a device
func (d *Device) Connect() (err *dbus.Error) {
	d.logger.Print("Device.Connect() not implemented")
	return err
}

// Disconnect a device
func (d *Device) Disconnect() (err *dbus.Error) {
	d.logger.Print("Device.Diconnect() not implemented")
	return nil
}

// Read data from a component
func (d *Device) Read(componentId string) (record *RecordObject, err *dbus.Error) {

	vm := otto.New()
	value, jserr := vm.Run(`return (new Date()).getTime();`)

	if jserr != nil {
		return nil, &dbus.Error{}
	}

	if val, err := value.ToString(); err == nil {

		d.logger.Printf("Result call for %s: %s\n", componentId, val)

		record = &RecordObject{}
		record.ComponentId = componentId
		record.DeviceId = d.Id
		record.Value = val
		record.Format = "int32"
		record.Unit = "whatever"

	}

	return nil, &dbus.Error{}
}

// Write data to a component
func (d *Device) Write(componentId string, record *RecordObject) (result bool, err *dbus.Error) {
	d.logger.Print("Device.Write() not implemented")
	return result, nil
}

// Subscribe for data updates to a component
func (d *Device) Subscribe(componentId string, params map[string]string) (result bool, err *dbus.Error) {
	d.logger.Print("Device.Subscribe() not implemented")
	return result, nil
}

// Unsubscribe from data updates of a component
func (d *Device) Unsubscribe(componentId string) (result bool, err *dbus.Error) {
	d.logger.Print("Device.Unsubscribe() not implemented")
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
