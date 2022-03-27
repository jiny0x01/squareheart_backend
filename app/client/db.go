package client

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jiny0x01/storylink_backend/ent"
	_ "github.com/lib/pq"
)

type DB struct {
	Client *ent.Client
	Ctx    context.Context
	Redis  *redis.Client
}

var db DB

func SetDB(d *DB) {
	db.Client = d.Client
	db.Ctx = d.Ctx
}

func GetDB() *DB {
	return &db
}

func init() {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	db.Redis = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err := db.Redis.Ping(db.Ctx).Result()
	if err != nil {
		panic(err)
	}
}
