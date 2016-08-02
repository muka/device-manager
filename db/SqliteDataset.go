package db

import (
	"database/sql"
	"fmt"
	"log"
	// "strconv"
	"strings"
	// use sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/muka/device-manager/util"
)

// NewSqliteDataSet instantiate a new sqlite dataset
func NewSqliteDataSet(tableName string, fields []DatasetField, filePath string) *SqliteDataset {
	logger, err := util.NewLogger("sqlite-dataset")
	util.CheckError(err)
	return &SqliteDataset{
		logger:    logger,
		tableName: tableName,
		fields:    fields,
		filePath:  filePath,
	}
}

// SqliteDataset a sqlite based dataset
type SqliteDataset struct {
	DataSet

	db      *sql.DB
	idField FieldValue
	logger  *log.Logger

	tableName string
	fields    []DatasetField
	filePath  string
}

// Open prepare for a database connection
func (d *SqliteDataset) Open() {

	if d.db != nil {
		return
	}

	d.logger.Printf("Open database at %s\n", d.filePath)
	db, err := sql.Open("sqlite3", d.filePath)
	util.CheckError(err)
	d.db = db

	fieldsList := ""
	for i, field := range d.fields {
		if i != 0 {
			fieldsList += "\n,"
		}
		if field.IsID {
			d.idField.Name = field.Name
		}
		fieldsList += fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Extras)
	}

	q := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` (\n%s)", d.tableName, fieldsList)
	d.logger.Printf("Ensured table exits %s:\n %s\n", d.tableName, q)

	_, err = d.db.Exec(q)
	util.CheckError(err)

	d.logger.Printf("Sqlite db ready\n")
}

// Close close the db connection
func (d *SqliteDataset) Close() {
	if d.db == nil {
		return
	}
	d.logger.Println("Closing connection")
	d.db.Close()
	d.db = nil
}

// Save a record
func (d *SqliteDataset) Save(fieldList []FieldValue) error {

	var idField FieldValue
	for _, field := range fieldList {
		if d.idField.Name == field.Name {
			idField = field
		}
	}

	d.logger.Printf("Saving record %s\n", idField.Value)

	fields := make([]string, len(fieldList))
	values := make([]string, len(fieldList))
	args := make([]interface{}, len(fieldList))

	for i, field := range fieldList {
		fields[i] = field.Name
		values[i] = "?" //"$" + strconv.Itoa(i+1)
		args[i] = field.Value
	}

	query := fmt.Sprintf("INSERT OR REPLACE INTO `%s` (%s) VALUES (%s)",
		d.tableName,
		strings.Join(fields, ","),
		strings.Join(values, ","))

	stmt, err := d.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(args...)

	return err
}

// Query execute a sql statment
func (d *SqliteDataset) Query(stmt string, args ...interface{}) (*sql.Rows, error) {
	d.Open()
	d.logger.Printf("Executing query\n %s -> %v\n", stmt, args)
	return d.db.Query(stmt, args...)
}

// GetBy get records for a field / value match
func (d *SqliteDataset) GetBy(key string, ids ...string) (*sql.Rows, error) {
	return d.Query("SELECT * FROM "+d.tableName+" WHERE ? IN (?)", key, ids)
}

// Find records in the db
func (d *SqliteDataset) Find(q *Query) (*sql.Rows, error) {

	var args = make([]interface{}, 0)

	var where string
	var orderby string

	// var p = 1
	if q != nil {

		if q.Criteria != nil {

			// p = len(q.Criteria)
			args = make([]interface{}, len(q.Criteria))
			var whereParts string
			for i, c := range q.Criteria {
				// where += c.Prefix + c.Field + c.Operation + "$" + strconv.Itoa(i+1) + c.Suffix
				whereParts += c.Prefix + c.Field + c.Operation + "?" + c.Suffix
				args[i] = c.Value
			}

			if where != "" {
				where = " WHERE " + whereParts
			}
		}

		if &q.OrderBy != nil {
			// stmt += " ORDER BY $" + strconv.Itoa(p) + " $" + strconv.Itoa(p+1)
			orderby += " ORDER BY ? ?"
			args = append(args, q.OrderBy.Field, q.OrderBy.Sort)
		}

	}

	var stmt = fmt.Sprintf("SELECT * FROM `%s` %s %s", d.tableName, where, orderby)
	return d.Query(stmt, args...)
}
