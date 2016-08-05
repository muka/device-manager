package service

import (
	"github.com/muka/device-manager/api"
	"github.com/muka/device-manager/util"
	"log"
)

var manager *Manager

// GetManager return the singleton Manager instance
func GetManager() *Manager {
	if manager == nil {
		manager = NewManager()
	}
	return manager
}

// NewManager creates a new Manager instance
func NewManager() *Manager {

	mlogger, err := util.NewLogger("service-manager")
	util.CheckError(err)

	s := Manager{}
	s.logger = mlogger
	s.instances = make(map[string]api.Service)

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

// Start create a service from a proxy and start it
func (s *Manager) Start(p api.Proxy) (*Service, error) {
	serviceInstance := NewService(p)
	err := s.StartService(serviceInstance)
	return serviceInstance, err
}

// Stop a running service
func (s *Manager) Stop(p api.Proxy) error {
	serviceInstance := s.instances[p.GetPath()]
	err := serviceInstance.Stop()
	if err != nil {
		return err
	}
	delete(s.instances, p.GetPath())
	return nil
}

// StartService start a service and add it to the managed list
func (s *Manager) StartService(service api.Service) error {

	var err error

	err = service.New()
	if err != nil {
		return err
	}

	err = service.Connect()
	if err != nil {
		return err
	}

	err = service.Export()
	if err != nil {
		return err
	}

	go service.Start()

	s.instances[service.GetPath()] = service
	s.logger.Printf("Added service for %s", service.GetPath())

	return nil
}
