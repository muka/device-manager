package dbus

import (
	"github.com/muka/device-manager/util"
	"log"
)

// New creates a new ServiceManager instance
func New() *ServiceManager {

	logger, err := util.NewLogger("service-manager")
	util.CheckError(err)

	return &ServiceManager{
		logger: logger,
	}
}

// ServiceManager handles many Service instances
type ServiceManager struct {
	instances []IService
	logger    *log.Logger
}

// Start a service and add it to the managed list
func (s *ServiceManager) Start(service IService) {

	var err error

	err = service.New()
	util.CheckError(err)

	err = service.Connect()
	util.CheckError(err)

	err = service.Export()
	util.CheckError(err)

	err = service.Start()
	util.CheckError(err)

}
