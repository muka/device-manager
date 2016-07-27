package main

import (
	"os"

	"github.com/godbus/dbus"
	"github.com/muka/device-manager/client"
	"github.com/muka/device-manager/config"
	"github.com/muka/device-manager/objects"
	"github.com/muka/device-manager/service"
	"github.com/muka/device-manager/util"
)

func init() {
	util.Logger()
	config.Get()
}

var startManager = func() {
	deviceManager := objects.NewDeviceManager()

	m := service.GetManager()
	_, err := m.Start(deviceManager)
	util.CheckError(err)
}

var createDevice = func() {
	dmClient := client.NewDeviceManager()

	dev := objects.DeviceDefinition{}

	dev.Id = "testId"
	dev.Name = "testName"
	dev.Description = "testDesc"

	dev.Path = dbus.ObjectPath("/iot/agile/my/TEST1")
	dev.Protocol = dbus.ObjectPath("/iot/agile/protocol/BLE")

	dev.Properties = make(map[string]string)
	dev.Properties["prop1"] = "propValue"

	dev.Streams = make([]objects.DeviceComponent, 1)

	dev.Streams[0] = objects.DeviceComponent{}
	dev.Streams[0].Id = "streamId"
	dev.Streams[0].Format = "streamFormat"
	dev.Streams[0].Unit = "streamUnit"
	dev.Streams[0].Properties = make(map[string]string)

	path, err := dmClient.Create(dev)
	util.CheckError(err)

	util.Logger().Printf("\nCreated device %s\n", path)
}

var readDevice = func(id string) {

	dmClient := client.NewDeviceManager()

	dev, err := dmClient.Read(id)
	util.CheckError(err)

	util.Logger().Printf("\nLoaded device %s\n", dev)
}

func main() {

	logger := util.Logger()

	if os.Args[1] == "" {
		logger.Print(`
Usage:
- "client" to run client
- "server" to run server
`)
		return
	}

	if os.Args[1] == "client" {
		// readDevice("test")
		createDevice()
	} else {
		startManager()
		select {}
	}
}
