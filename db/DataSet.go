package db

import (
	"database/sql"
)

// DatasetField a field of a dataset
type DatasetField struct {
	Name   string
	Type   string
	Extras string
	IsID   bool
}

// FieldValue a field mapped to a value
type FieldValue struct {
	Name  string
	Value interface{}
}

//Sort sorting direction, ASC or DESC
type Sort string

//OrderBy sorting direction, ASC or DESC
type OrderBy struct {
	Field string
	Sort  Sort
}

const (
	//SortASC sort ascending
	SortASC Sort = "ASC"
	//SortDESC sort descending
	SortDESC Sort = "DESC"
)

// Criteria a single quert criteria eg. [name] [=] ["something"]
type Criteria struct {
	Prefix    string
	Field     string
	Operation string
	Value     string
	Suffix    string
}

//Limit the limit and offset of the query
type Limit struct {
	//Offset nr of record to skip
	Offset int
	// Size length of the set to return
	Size int
}

// Query abstracts a query to the db
type Query struct {
	// Criteria list of WHERE parameters
	Criteria []Criteria
	// OrderBy the order of the query
	OrderBy OrderBy
	Limit   Limit
}

// DataSet abstracts a storage engine
type DataSet interface {
	Open()
	Close()
	GetBy(key string, ids ...string) (*sql.Rows, error)
	Find(q *Query) (*sql.Rows, error)
	Query(stmt string, args ...interface{}) (*sql.Rows, error)
	Save(fieldList []FieldValue) error
}
