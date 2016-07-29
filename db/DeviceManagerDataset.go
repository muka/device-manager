package db

import (
	"database/sql"
	"github.com/muka/device-manager/api"
	"github.com/muka/device-manager/util"
	// use sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

const (
	defaultFilePath  = "./data.db"
	defaultTableName = "Devices"
)

// DeviceManagerDataset a sqlite based dataset
type DeviceManagerDataset struct {
	db *sql.DB
}

// Setup a database connection
func (d *DeviceManagerDataset) Setup() {

	sqlStmt := `
  create table ` + defaultTableName + ` (
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
	_, err := d.db.Exec(sqlStmt)
	util.CheckError(err)

}

// Connect return a connection to the database
func (d *DeviceManagerDataset) Connect() {
	db, err := sql.Open("sqlite3", defaultFilePath)
	util.CheckError(err)
	d.db = db
}

// Disconnect close the db connection
func (d *DeviceManagerDataset) Disconnect() {
	d.db.Close()
}

// GetBy get records for a field / value match
func (d *DeviceManagerDataset) GetBy(key string, id string) *sql.Result {
	result, err := d.db.Exec("SELECT * FROM "+defaultTableName+" WHERE $1 = $2", key, id)
	util.CheckError(err)
	return &result
}

// Find records in the db
func (d *DeviceManagerDataset) Find(q api.Query) *sql.Result {
	result, err := d.db.Exec("SELECT * FROM "+defaultTableName+" WHERE Id = $1 ORDER BY $2 $3", q.Criteria, q.OrderBy, q.Sort)
	util.CheckError(err)
	return &result
}
