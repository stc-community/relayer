package sqlite3

import (
	"github.com/jmoiron/sqlx"
)

type SQLite3 struct {
	*sqlx.DB
	DatabaseURL string
}
