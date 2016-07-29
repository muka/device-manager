package api

import (
	"database/sql"
)

// Query abstracts a query to the db
type Query struct {
	Criteria map[string]string
	OrderBy  string
	Sort     string
}

// DataSet abstracts a storage engine
type DataSet interface {
	Setup()
	Connect()
	Disconnect()
	GetBy(key string, id string) *sql.Result
	Find(q Query) *sql.Result
}
