package service

import (
	"github.com/muka/device-manager/api"
	"github.com/muka/device-manager/util"
	"log"
)

// NewManager creates a new Manager instance
func NewManager() *Manager {

	mlogger, err := util.NewLogger("service-manager")
	util.CheckError(err)

	s := Manager{}
	s.logger = mlogger

	return &s
}

// Manager handles Service instances
type Manager struct {
	instances map[string]api.Service
	logger    *log.Logger
}

// GetService return an instance of a service
func (s *Manager) GetService(path string) api.Service {
	return s.instances[path]
}

// Start a service and add it to the managed list
func (s *Manager) Start(service api.Service) {

	var err error

	err = service.New()
	util.CheckError(err)

	err = service.Connect()
	util.CheckError(err)

	err = service.Export()
	util.CheckError(err)

	err = service.Start()
	util.CheckError(err)

	s.instances[service.GetPath()] = service
	s.logger.Printf("Added service for %s", service.GetPath())
}
