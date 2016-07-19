package main

import (
	"github.com/muka/device-manager/config"
	"github.com/muka/device-manager/dbus"
	"github.com/muka/device-manager/objects"
	"github.com/muka/device-manager/util"
)

var logger = util.Logger()

func main() {

	config := config.Get()
	logger.Printf("Config %v\n", config)

	daemon := dbus.Daemon{
		Path:   "/iot/agile/DeviceManager",
		Iface:  "iot.agile.DeviceManager",
		Object: &objects.DeviceManager{},
	}

	daemon.Start()
}
