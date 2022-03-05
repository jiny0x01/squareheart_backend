package client

import (
	"context"

	"github.com/jiny0x01/storylink_backend/ent"
	_ "github.com/lib/pq"
)

type DB struct {
	Client *ent.Client
	Ctx    context.Context
}

var db DB

func SetDB(d *DB) {
	db.Client = d.Client
	db.Ctx = d.Ctx
}

func GetDB() *DB {
	return &db
}
