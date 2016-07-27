package objects

// RecordObject A record of a device sensor
type RecordObject struct {
	Format      string
	Value       string
	Unit        string
	ComponentId string
	DeviceId    string
	LastUpdate  int32
}
