package validata

import "database/sql"

type Messages []string

type Math interface {
	Analyze(data []interface{}) error
}

type Notify interface {
	Notify(Messages) error
}

type DB interface {
	Open(driverName, dataSourceName string) error
	Query(args ...interface{}) ([]sql.Rows, error)
	Close() error
}
