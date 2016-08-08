package main

import (
	"github.com/godbus/dbus"
	"github.com/muka/device-manager/client"
	"github.com/muka/device-manager/objects"
	"github.com/muka/device-manager/service"
	"os"
	"testing"
)

const defaultDeviceID = "TEST1"
const defaultDeviceStreamID = "temperature"

var testDevice objects.DeviceDefinition

var deviceManagerService *service.Service

func TestMain(t *testing.M) {

	var err error
	deviceManagerService, err = StarDeviceManager()

	if err != nil {
		os.Exit(1)
	}

	t.Run()
}

func newDeviceOverview() objects.DeviceDefinition {

	dev := objects.DeviceDefinition{}

	dev.Id = defaultDeviceID
	dev.Name = "TEST device"
	dev.Description = "A device created during test"

	dev.Path = dbus.ObjectPath("/iot/agile/my/" + dev.Id)
	dev.Protocol = dbus.ObjectPath("/iot/agile/protocol/BLE")

	dev.Properties = make(map[string]string)
	dev.Properties["prop1"] = "propValue"

	dev.Streams = make([]objects.DeviceComponent, 1)

	dev.Streams[0] = objects.DeviceComponent{}
	dev.Streams[0].Id = defaultDeviceStreamID
	dev.Streams[0].Format = "double"
	dev.Streams[0].Unit = "celsius"
	dev.Streams[0].Properties = make(map[string]string)

	return dev
}

func TestCreateDevice(t *testing.T) {
	deviceManagerClient := client.NewDeviceManager()

	device, err := deviceManagerClient.Create(newDeviceOverview())
	if err != nil {
		t.Fail()
	}

	testDevice = device
	t.Logf("Created device with ID %s", testDevice.Id)
}

func TestReadDevice(t *testing.T) {

	deviceManagerClient := client.NewDeviceManager()
	dev, err := deviceManagerClient.Read(testDevice.Id)

	if err != nil {
		t.Fail()
	}

	if len(dev.Streams) == 1 {
		if dev.Streams[0].Id == defaultDeviceStreamID {
			return
		}
	}
	t.Fail()
}

func TestFindDevice(t *testing.T) {

	deviceManagerClient := client.NewDeviceManager()
	dev, err := deviceManagerClient.Find()

	if err != nil {
		t.Fail()
	}

	if len(dev.Streams) == 1 {
		if dev.Streams[0].Id == defaultDeviceStreamID {
			return
		}
	}
	t.Fail()
}
