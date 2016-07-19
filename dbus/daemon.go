package dbus

import (
	"fmt"
	"regexp"

	"github.com/muka/device-manager/config"
	"github.com/muka/device-manager/fetch"
	"github.com/muka/device-manager/objects"
	"github.com/muka/device-manager/util"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

var logger = util.Logger()

func check(err error) {
	if err != nil {
		logger.Fatal(err)
		panic(err)
	}
}

// Daemon provide deamon instance informations
type Daemon struct {
	Path    dbus.ObjectPath
	Iface   string
	Object  objects.Proxy
	XMLPath string
	XML     string
}

// PrepareXML format the XML to be passed to dbus
func (d Daemon) PrepareXML() {

	xml, err := fetch.GetContent(d.XMLPath)
	check(err)

	reg := regexp.MustCompile(`(?s).*<node`)
	xml = reg.ReplaceAllString(xml, "<node")

	reg = regexp.MustCompile(`(<\/node>)`)
	xml = reg.ReplaceAllString(xml,
		introspect.IntrospectDataString+"</node>")

	logger.Printf("DBus XML:\n\n %s \n\n", xml)

	d.XML = xml

}

// Start start a new daemon
func (d Daemon) Start() {

	config := config.Get()

	specPath := config.APISpecPath
	absSpecPath, err := fetch.ResolvePath(specPath)
	check(err)

	d.XMLPath = absSpecPath + "/" + d.Iface + ".xml"

	d.PrepareXML()

	conn, err := dbus.SessionBus()
	check(err)

	reply, err := conn.RequestName(d.Iface,
		dbus.NameFlagDoNotQueue)
	check(err)

	if reply != dbus.RequestNameReplyPrimaryOwner {
		msg := "name " + d.Iface + " already taken"
		logger.Fatalln(msg)
		panic(msg)
	}

	obj := d.Object

	conn.Export(obj, d.Path, d.Iface)
	conn.Export(introspect.Introspectable(d.XML),
		d.Path,
		"org.freedesktop.DBus.Introspectable")

	fmt.Printf("Listening on %s %s ...", d.Iface, d.Path)

	select {}
}
