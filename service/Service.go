package service

import (
	"errors"
	"log"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
	"github.com/godbus/dbus/prop"

	"github.com/muka/device-manager/api"
	"github.com/muka/device-manager/util"
)

// NewService instantiate a new service
func NewService(p api.Proxy) *Service {

	s := Service{}

	s.SetObject(&p)

	s.SetPath(p.GetPath())
	s.SetInterface(p.GetInterface())

	return &s
}

// Service provide deamon instance informations
type Service struct {
	api.Service

	object *api.Proxy
	path   string
	iface  string

	conn     *dbus.Conn
	dbusPath dbus.ObjectPath
	logger   *log.Logger
}

//GetPath return the service path
func (d *Service) GetPath() string {
	return d.path
}

//SetPath set the service path
func (d *Service) SetPath(path string) {
	d.path = path
}

//GetInterface return the service interface
func (d *Service) GetInterface() string {
	return d.iface
}

//SetInterface set the service interface
func (d *Service) SetInterface(iface string) {
	d.iface = iface
}

//SetObject set the object proxy
func (d *Service) SetObject(p *api.Proxy) {
	d.object = p
}

//GetObject get the object proxy
func (d *Service) GetObject() *api.Proxy {
	return d.object
}

// New configure a daemon instance
func (d *Service) New() error {

	dlogger, err := util.NewLogger(d.GetInterface())
	if err != nil {
		return err
	}

	d.logger = dlogger
	(*d.GetObject()).SetLogger(dlogger)

	d.dbusPath = dbus.ObjectPath(d.GetPath())

	d.logger.Println("Setup completed")
	return nil
}

// Connect connect to a bus
func (d *Service) Connect() error {

	conn, err := dbus.SessionBus()

	if err != nil {
		return err
	}

	d.conn = conn

	reply, err := d.conn.RequestName(d.GetInterface(), dbus.NameFlagDoNotQueue)
	if err != nil {
		return err
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		return errors.New("Name " + d.GetInterface() + " already taken")
	}

	d.logger.Println("Connect completed")
	return nil
}

// Export the daemon node & intefaces
func (d *Service) Export() error {

	d.logger.Printf("Export to Dbus %s\n", d.GetPath())

	// Expose the object path
	d.conn.Export(d.GetObject(),
		d.dbusPath,
		d.GetInterface())

	// Register properties
	propsSpec := (*d.GetObject()).GetProperties()
	props := prop.New(d.conn, d.dbusPath, propsSpec)
	root := &introspect.Node{
		Name: d.GetPath(),
		Interfaces: []introspect.Interface{
			introspect.IntrospectData,
			prop.IntrospectData,
			{
				Name:       d.GetInterface(),
				Methods:    introspect.Methods(d.GetObject()),
				Properties: props.Introspection(d.GetInterface()),
			},
		},
	}

	// Export Introspectable
	d.conn.Export(introspect.NewIntrospectable(root),
		d.dbusPath,
		"org.freedesktop.DBus.Introspectable")

	return nil
}

// Start start a new daemon
func (d *Service) Start() error {

	// config := config.Get()
	d.logger.Printf("Started listening on %s %s ...\n",
		d.GetInterface(), d.GetPath())

	select {}
}

// Stop stop the daemon
func (d *Service) Stop() error {
	if d.conn != nil {
		d.conn.Export(nil, d.dbusPath, d.GetInterface())
		d.conn.ReleaseName(d.GetInterface())
		d.logger.Printf("Service %s stopped\n", d.GetPath())
	} else {
		d.logger.Println("Service already stopped")
	}
	return nil
}
