package xmobile

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"go.olapie.com/x/xsqlite"
)

type SQLite3DB struct {
	db    *sql.DB
	error string
}

func (s *SQLite3DB) DB() *sql.DB {
	return s.db
}

func (s *SQLite3DB) Error() string {
	return s.error
}

func (s *SQLite3DB) Close() {
	if s.db != nil {
		err := s.db.Close()
		if err != nil {
			s.error = err.Error()
		}
	}
}

func NewSQLite3DB(filename string) *SQLite3DB {
	db, err := xsqlite.Open(filename)
	if err != nil {
		return &SQLite3DB{
			error: err.Error(),
		}
	}
	return &SQLite3DB{
		db: db,
	}
}
