package api

import (
	"github.com/godbus/dbus/prop"
	"log"
)

//Proxy object that can be exposed as DBus service
type Proxy interface {
	GetPath() string
	SetPath(s string)

	GetInterface() string
	SetInterface(s string)

	SetLogger(logger *log.Logger)
	GetLogger() *log.Logger

	GetProperties() map[string]map[string]*prop.Prop
}
