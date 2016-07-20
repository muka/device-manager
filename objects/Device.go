package objects

import (
	"log"
)

// Device a device implementation
type Device struct {
	BaseObject
	Logger *log.Logger
}
