package api

import (
	"github.com/godbus/dbus/prop"
	"log"
)

//Proxy object that can be exposed as DBus service
type Proxy interface {
	GetPath() string
	GetInterface() string

	GetProperties() map[string]map[string]*prop.Prop
	SetLogger(logger *log.Logger)
}
