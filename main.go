package main

import (
	"github.com/muka/device-manager/config"
	"github.com/muka/device-manager/objects"
	"github.com/muka/device-manager/service"
	"github.com/muka/device-manager/util"
)

func init() {
	util.Logger()
	config.Get()
}

func main() {

	deviceManager := objects.NewDeviceManager()

	m := service.NewManager()
	m.Start(service.NewService(deviceManager))
}
