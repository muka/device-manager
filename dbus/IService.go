package dbus

// IService generic daemon interface
type IService interface {
	New() error
	Connect() error
	Export() error
	Start() error
	Stop() error
}
