package dbus

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/muka/device-manager/config"
	"github.com/muka/device-manager/fetch"
	"github.com/muka/device-manager/objects"
	"github.com/muka/device-manager/util"

	"github.com/godbus/dbus"
	"github.com/godbus/dbus/introspect"
)

var logger = util.Logger()

const (
	appendString = "<!-- @EndOfNode -->"
	appendReg    = `\<\!\-\- \@EndOfNode \-\-\>`
)

func check(err error) {
	if err != nil {
		logger.Fatal(err)
		panic(err)
	}
}

// Daemon provide deamon instance informations
type Daemon struct {
	Path    string
	Iface   string
	Object  objects.IProxy
	XMLPath string
	XML     string
	conn    *dbus.Conn
}

// LoadXML format the XML to be passed to dbus
func (d *Daemon) LoadXML() {

	xml, err := fetch.GetContent(d.XMLPath)
	check(err)

	reg := regexp.MustCompile(`(?s).*<node`)
	xml = reg.ReplaceAllString(xml, "<node")

	// add append reference
	reg = regexp.MustCompile(`(<\/node>)`)
	xml = reg.ReplaceAllString(xml, appendString+"\n</node>")

	d.XML = xml

	// add Introspection string
	d.appendToXMLNode(introspect.IntrospectDataString)

	// exports the tree of introspectable
	d.addIntrospectableChildren()
}

func (d *Daemon) appendToXMLNode(str string) {
	reg := regexp.MustCompile(appendReg)
	d.XML = reg.ReplaceAllString(d.XML, str+appendString)
}

func (d *Daemon) addIntrospectableChildren() {

	parts := strings.Split(d.Path, "/")
	partialPath := ""
	xml := ""

	partsLen := len(parts)
	for i := 0; i < partsLen; i++ {

		if i+1 == partsLen {
			break
		}

		path := parts[i]
		switch i {
		case 0:
			partialPath = "/"
		case 1:
			partialPath += path
		default:
			partialPath += "/" + path
		}
		xml += `<node name="` + partialPath + `" />`
	}

	if xml != "" {
		d.appendToXMLNode(xml)
	}
}

// Start start a new daemon
func (d *Daemon) Start() {

	config := config.Get()

	specPath := config.APISpecPath
	absSpecPath, err := fetch.ResolvePath(specPath)
	check(err)

	d.XMLPath = absSpecPath + "/" + d.Iface + ".xml"
	d.LoadXML()

	conn, err := dbus.SessionBus()
	check(err)

	d.conn = conn

	reply, err := conn.RequestName(d.Iface,
		dbus.NameFlagDoNotQueue)
	check(err)

	if reply != dbus.RequestNameReplyPrimaryOwner {
		msg := "name " + d.Iface + " already taken"
		logger.Fatalln(msg)
		panic(msg)
	}

	obj := d.Object
	objPath := dbus.ObjectPath(d.Path)

	logger.Printf("DBus XML:\n\n %s \n\n", d.XML)

	err = obj.DBusSetup(conn)
	check(err)

	conn.Export(introspect.Introspectable(d.XML),
		objPath,
		"org.freedesktop.DBus.Introspectable")

	conn.Export(obj, objPath, d.Iface)

	fmt.Printf("Listening on %s %s ...", d.Iface, d.Path)

	select {}
}
