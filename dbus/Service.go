package dbus

import (
	"errors"
	"github.com/godbus/dbus"
	"github.com/muka/device-manager/objects"
	"github.com/muka/device-manager/util"
	"log"
)

// Service provide deamon instance informations
type Service struct {
	Object    objects.BaseObject
	Path      string
	Interface string

	conn     *dbus.Conn
	dbusPath dbus.ObjectPath
	logger   *log.Logger
}

// New configure a daemon instance
func (d *Service) New() error {

	dlogger, err := util.NewLogger(d.Interface)
	if err != nil {
		return err
	}

	d.logger = dlogger
	d.Object.Logger = dlogger

	d.dbusPath = dbus.ObjectPath(d.Path)

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

	reply, err := d.conn.RequestName(d.Interface, dbus.NameFlagDoNotQueue)
	if err != nil {
		return err
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		return errors.New("Name " + d.Interface + " already taken")
	}

	d.logger.Println("Connect completed")
	return nil
}

// Start start a new daemon
func (d *Service) Start() error {

	// config := config.Get()
	d.logger.Printf("Started listening on %s %s ...\n",
		d.Interface, d.Path)

	select {}
}

// Stop stop the daemon
func (d *Service) Stop() error {
	if d.conn != nil {
		d.conn.Export(nil, d.dbusPath, d.Interface)
		d.conn.ReleaseName(d.Interface)
		d.logger.Println("Service stopped")
	} else {
		d.logger.Println("Service already stopped")
	}
	return nil
}
