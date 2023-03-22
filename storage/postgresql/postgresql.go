package postgresql

import (
	"github.com/jmoiron/sqlx"
)

type Postgres struct {
	*sqlx.DB
	DatabaseURL string
}
