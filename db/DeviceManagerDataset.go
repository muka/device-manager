package db

import (
	"database/sql"
	"github.com/muka/device-manager/api"
	"github.com/muka/device-manager/util"
	// use sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

const (
	// DatabasePath path to the database file
	DatabasePath = "./data.db"
	// TableName Name of the table containing the devices
	TableName = "Devices"

	createTableStmt = `
  create table ` + TableName + ` (
    id integer not null primary key,
		Id          text,
		Name        text,
		Description text,
		Path        text,
		Protocol    text,
		Properties  text,
		Streams     text
  );
  `
)

// DeviceManagerDataset a sqlite based dataset
type DeviceManagerDataset struct {
	db *sql.DB
}

// Open setup internally for a database connection
func (d *DeviceManagerDataset) Open() {

	db, err := sql.Open("sqlite3", DatabasePath)
	util.CheckError(err)
	d.db = db

	_, err = d.db.Exec(createTableStmt)
	util.CheckError(err)

}

// Connect return a connection to the database
func (d *DeviceManagerDataset) Connect() {

}

// Disconnect close the db connection
func (d *DeviceManagerDataset) Disconnect() {
	d.db.Close()
}

// GetBy get records for a field / value match
func (d *DeviceManagerDataset) GetBy(key string, id string) *sql.Result {
	result, err := d.db.Exec("SELECT * FROM "+TableName+" WHERE $1 = $2", key, id)
	util.CheckError(err)
	return &result
}

// Find records in the db
func (d *DeviceManagerDataset) Find(q api.Query) *sql.Result {
	var where = ""
	var args = make([]interface{}, len(q.Criteria))
	for i, c := range q.Criteria {
		p := string(i + 1)
		where += c.Prefix + c.Field + c.Operation + "$" + p + c.Suffix
		args = append(args, c.Value)
	}
	result, err := d.db.Exec("SELECT * FROM "+TableName+" WHERE "+where+" ORDER BY $2 $3", q.Criteria, q.OrderBy.Field, q.OrderBy.Sort)
	util.CheckError(err)
	return &result
}
