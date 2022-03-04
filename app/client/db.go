package client

import (
	"github.com/jiny0x01/storylink_backend/ent"
	_ "github.com/lib/pq"
)

type DB struct {
	*ent.Client
}

var db DB

func SetDB(d *DB) {
	db.Client = d.Client
}

func GetDB() *DB {
	return &db
}
