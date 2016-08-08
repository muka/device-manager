package objects

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/prop"

	"github.com/muka/device-manager/api"
	"github.com/muka/device-manager/db"
	"github.com/muka/device-manager/service"
	"github.com/muka/device-manager/util"
)

// DeviceManagerDatabasePath path to the database file
var deviceManagerDatabasePath = "./data/devices.db"

// DeviceManagerTableName Name of the table containing the devices
var deviceManagerTableName = "Devices"

var deviceManagerFields = []db.DatasetField{
	{"Id", "text", "unique", true},
	{Name: "Name", Type: "text"},
	{Name: "Description", Type: "text"},
	{Name: "Path", Type: "text"},
	{Name: "Protocol", Type: "text"},
	{Name: "Properties", Type: "text"},
	{Name: "Streams", Type: "text"},
}

// NewDeviceManager initialize a new DeviceManager object
func NewDeviceManager() *DeviceManager {
	d := DeviceManager{}

	dmlogger, err := util.NewLogger("DeviceManager")
	util.CheckError(err)

	d.logger = dmlogger

	d.Devices = make([]dbus.ObjectPath, 0)
	d.services = make(map[string]*service.Service)
	d.devices = make(map[string]*DeviceDefinition)
	d.path = DeviceManagerPath
	d.iface = DeviceManagerInterface
	d.dataset = db.NewSqliteDataSet(deviceManagerTableName, deviceManagerFields, deviceManagerDatabasePath)

	d.restoreDevices()

	return &d
}

// DeviceManager manages devices in the gateway
type DeviceManager struct {
	api.Proxy

	Devices  []dbus.ObjectPath
	services map[string]*service.Service
	devices  map[string]*DeviceDefinition

	path    string
	iface   string
	logger  *log.Logger
	dataset db.DataSet
}

// GetPath return object path
func (d *DeviceManager) GetPath() string {
	return d.path
}

// SetPath set object path
func (d *DeviceManager) SetPath(s string) {
	d.path = s
}

// GetInterface return interface
func (d *DeviceManager) GetInterface() string {
	return d.iface
}

// SetInterface return interface
func (d *DeviceManager) SetInterface(s string) {
	d.iface = s
}

//SetLogger set default logger
func (d *DeviceManager) SetLogger(logger *log.Logger) {
	d.logger = logger
}

//GetLogger return default logger
func (d *DeviceManager) GetLogger() *log.Logger {
	return d.logger
}

//GetProperties return properties
func (d *DeviceManager) GetProperties() map[string]map[string]*prop.Prop {
	return map[string]map[string]*prop.Prop{
		d.GetInterface(): {
			"Devices": {
				Value:    d.Devices,
				Writable: false,
				Emit:     prop.EmitTrue,
				Callback: func(c *prop.Change) *dbus.Error {
					d.logger.Printf("Changed value %s=%v on %s", c.Name, c.Value, c.Iface)
					return nil
				},
			},
		},
	}
}

//parseDeviceRow parse a row to DeviceDefinition
func (d *DeviceManager) parseDeviceRow(rows *sql.Rows) (*DeviceDefinition, error) {
	dev := new(DeviceDefinition)
	var strStreams, strProperties, strPath, strProtocol string
	err := rows.Scan(
		&dev.Id,
		&dev.Name,
		&dev.Description,
		&strPath,
		&strProtocol,
		&strProperties,
		&strStreams,
	)
	if err != nil {
		return dev, err
	}

	dev.Protocol = dbus.ObjectPath(strProtocol)
	dev.Path = dbus.ObjectPath(strPath)

	json.Unmarshal([]byte(strProperties), dev.Properties)
	json.Unmarshal([]byte(strProperties), dev.Streams)

	return dev, nil
}

//restoreDevices reinitialize DBus instances of stored devices
func (d *DeviceManager) restoreDevices() {

	d.logger.Println("Restoring previous device instances")

	rows, err := d.dataset.Find(nil)
	util.CheckError(err)

	defer rows.Close()
	for rows.Next() {
		dev, err := d.parseDeviceRow(rows)
		util.CheckError(err)
		d.logger.Printf("Loading device %s (%s)", dev.Name, dev.Id)
		d.startDeviceInstance(*dev)
	}

}

//startDeviceInstance start a device instance daemon service
func (d *DeviceManager) startDeviceInstance(dev DeviceDefinition) error {

	device := NewDevice(dev)

	service, mErr := service.GetManager().Start(device)
	if mErr != nil {
		msg := fmt.Sprintf("Cannot start Device service: %s\n", mErr.Error())
		d.logger.Fatalf(msg)
		return errors.New(msg)
	}

	dev.Path = dbus.ObjectPath(device.GetPath())

	d.Devices = append(d.Devices, dev.Path)
	d.devices[dev.Id] = &dev
	d.services[dev.Id] = service

	return nil
}

//loadDevice load from database a device definition
func (d *DeviceManager) loadDevice(id string) (*DeviceDefinition, error) {
	rows, err := d.dataset.GetBy("Id", id)
	if err != nil {
		d.logger.Fatalf("Cannot load device record %s", id)
		return nil, err
	}

	var dev DeviceDefinition
	defer rows.Close()
	for rows.Next() {
		dev, err := d.parseDeviceRow(rows)
		util.CheckError(err)
		d.logger.Printf("Loaded record for device %s (%s)", dev.Name, dev.Id)
		break
	}

	return &dev, nil
}

//saveDevice save to database a device definition
func (d *DeviceManager) saveDevice(dev DeviceDefinition) error {

	jsonProperties, err := json.Marshal(dev.Properties)
	if err != nil {
		return err
	}

	jsonStreams, err := json.Marshal(dev.Streams)
	if err != nil {
		return err
	}

	err = d.dataset.Save(
		[]db.FieldValue{
			{Name: "Id", Value: dev.Id},
			{Name: "Name", Value: dev.Name},
			{Name: "Description", Value: dev.Description},
			{Name: "Path", Value: string(dev.Path)},
			{Name: "Protocol", Value: string(dev.Protocol)},
			{Name: "Properties", Value: string(jsonProperties)},
			{Name: "Streams", Value: string(jsonStreams)},
		},
	)

	if err != nil {
		d.logger.Printf("Error on save: %v\n", err)
		return err
	}

	return nil
}

// -----
// Dbus API implementation

// Find search for devices
func (d *DeviceManager) Find(q *BaseQuery) (devices []dbus.ObjectPath, err *dbus.Error) {

	var query = db.Query{}
	if q.Criteria != nil {
		query.Criteria = make([]db.Criteria, len(q.Criteria))
		var i = 0
		for key, val := range q.Criteria {

			var op = "="
			var value = val

			if strings.Contains(value, "*") {
				op = "LIKE"
				value = strings.Replace(value, "*", "%", 0)
			}

			query.Criteria[i] = db.Criteria{
				Prefix:    "",
				Field:     key,
				Operation: op,
				Value:     value,
				Suffix:    "",
			}

			i++
		}
	}

	query.Limit = db.Limit{}
	if q.Offset > 0 {
		query.Limit.Offset = int(q.Offset)
	}
	if q.Limit > 0 {
		query.Limit.Size = int(q.Limit)
	}

	for k, v := range q.OrderBy {
		s := db.SortDESC
		if v == "ASC" {
			s = db.SortASC
		}
		query.OrderBy = db.OrderBy{Field: k, Sort: s}
		break
	}

	rows, err1 := d.dataset.Find(&query)
	util.CheckError(err1)

	defer rows.Close()
	var i = 0
	var devs = make([]dbus.ObjectPath, 0)
	for rows.Next() {
		dev, err := d.parseDeviceRow(rows)
		util.CheckError(err)
		devs = append(devs, dev.Path)
		i++
	}

	return devs, err
}

// Create add a device
func (d *DeviceManager) Create(dev DeviceDefinition) (DeviceDefinition, *dbus.Error) {

	var err error

	id := util.GenerateID()
	d.logger.Printf("Create new device %s\n", id)
	d.logger.Printf("Data:\n %v\n", dev)
	dev.Id = id

	d.logger.Printf("Save record for device %s\n", dev.Id)
	err = d.saveDevice(dev)
	if err != nil {
		d.logger.Fatalf("Error saving record for device %s\n", dev.Id)
		return DeviceDefinition{}, &dbus.Error{}
	}

	err = d.startDeviceInstance(dev)
	if err != nil {
		d.logger.Fatalf("Error starting service for device %s\n", dev.Id)
		return DeviceDefinition{}, &dbus.Error{}
	}

	d.logger.Printf("Created new device %s\n", dev.Id)
	return dev, nil
}

// Read a device definition
func (d *DeviceManager) Read(id string) (dev DeviceDefinition, err *dbus.Error) {

	if d.devices[id] != nil {
		dev = *d.devices[id]
		d.logger.Printf("Read %s: \n%v\n", id, dev)
	} else {
		d.logger.Printf("Device %s: Not Found : \n", id)
		err = new(dbus.Error)
	}

	return dev, err
}

// Update a device definition
func (d *DeviceManager) Update(id string, dev DeviceDefinition) (res bool, err *dbus.Error) {

	// Check if device exists
	if localDev, exists := d.devices[id]; exists {

		localDev.Name = dev.Name
		localDev.Description = dev.Description
		localDev.Path = dev.Path
		localDev.Protocol = dev.Protocol
		localDev.Properties = dev.Properties
		localDev.Streams = dev.Streams

		d.saveDevice(*localDev)

		d.logger.Printf("Device %s updated\n", id)
		return true, err
	}

	d.logger.Fatalf("Device %s not found\n", id)
	return false, new(dbus.Error)
}

// Delete a device definition
func (d *DeviceManager) Delete(id string) (res bool, err *dbus.Error) {

	d.logger.Printf("Removing device %s\n", id)

	if _, exists := d.devices[id]; !exists {
		d.logger.Fatalf("Device %s not found\n", id)
		return false, new(dbus.Error)
	}

	devService := d.services[id]

	mErr := devService.Stop()
	if mErr != nil {
		d.logger.Fatalf("Cannot stop service for device %s\n", id)
		return false, new(dbus.Error)
	}

	err1 := d.dataset.Delete(db.FieldValue{Name: "Id", Value: id})
	if err1 == nil {
		d.logger.Fatalf("Cannot stop service for device %s\n", id)
		return false, new(dbus.Error)
	}

	delete(d.devices, id)
	delete(d.services, id)

	d.logger.Printf("Removed device %s\n", id)
	return true, err
}

// Batch exec batch ops
// func (d *DeviceManager) Batch(operation string, ops map[string]string) (res bool, err *dbus.Error) {
// 	res = true
// 	return res, err
// }

// BaseQuery base query for devices record
type BaseQuery struct {
	Criteria map[string]string
	OrderBy  map[string]string
	Limit    int32
	Offset   int32
}

// DeviceComponent A device component
type DeviceComponent struct {
	Id         string
	Unit       string
	Format     string
	Properties map[string]string
}

// DeviceDefinition A device details list
type DeviceDefinition struct {
	Id          string
	Name        string
	Description string
	Path        dbus.ObjectPath
	Protocol    dbus.ObjectPath
	Properties  map[string]string
	Streams     []DeviceComponent
}
