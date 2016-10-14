package main

import (
	"github.com/muka/device-manager/client"
	"github.com/muka/device-manager/config"
	"github.com/muka/device-manager/objects"
	"github.com/muka/device-manager/service"
	"github.com/muka/device-manager/util"
	"os"
)

func init() {
	util.Logger()
	config.Get()
}

// StarDeviceManager start an instance of the device manager
func StarDeviceManager() (*service.Service, error) {
	deviceManager := objects.NewDeviceManager()
	m := service.GetManager()
	return m.Start(deviceManager)
}

// ReadDevice read a device by Id
var ReadDevice = func(id string) {

	dmClient := client.NewDeviceManager()

	dev, err := dmClient.Read(id)
	util.CheckError(err)

	util.Logger().Printf("\nLoaded device %s\n", dev)
}

func main() {

	logger := util.Logger()

	if len(os.Args) == 1 {
		logger.Print(`
	Usage:
	- "client" to run client
	- "server" to run server
	`)
		return
	}

	if os.Args[1] == "client" {
		ReadDevice("test")
	} else {
		StarDeviceManager()
		select {}
	}
}
