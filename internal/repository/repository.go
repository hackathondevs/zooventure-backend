package repository

import "github.com/jmoiron/sqlx"

type txCompat interface {
	Commit() error
	Rollback() error
	Ext() sqlx.ExtContext
}
