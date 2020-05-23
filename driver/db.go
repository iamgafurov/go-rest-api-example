package db


import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

type DB struct {
	*sql.DB
}

func ConnectDB(dbSource, dbType string) (*DB, error) {

	d, error := sql.Open(dbType, dbSource)

	if error != nil {
		return nil, error
	}

	if err := d.Ping(); err != nil {
		return nil, err
	}

	return &DB{d}, nil

}

