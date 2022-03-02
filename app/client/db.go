package client

import (
	"storylink_backend/ent"

	_ "github.com/lib/pq"
)

type DB ent.Client

var db *DB

func SetDB(d *DB) {
	db = d
}

func GetDB() *DB {
	return db
}
