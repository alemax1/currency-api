package postgres

import "database/sql"

type DBExecutor struct {
	db *sql.DB
}

func NewExecutor(db *sql.DB) *DBExecutor {
	return &DBExecutor{
		db: db,
	}
}
