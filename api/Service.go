package api

// Service generic daemon interface
type Service interface {
	New() error
	Connect() error
	Export() error
	Start() error
	Stop() error

	GetObject() *Proxy
	SetObject(p *Proxy)

	GetPath() string
	SetPath(path string)

	GetInterface() string
	SetInterface(iface string)
}
