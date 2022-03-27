package client

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/jiny0x01/storylink_backend/ent"
	"github.com/joho/godotenv"
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
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error db loading .env file")
		return
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s TimeZone=Asia/Seoul",
		os.Getenv("HOST"),
		os.Getenv("PORT"),
		os.Getenv("STORYLINK_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("STORYLINK_PW"),
		os.Getenv("SSL_MODE"))
	log.Printf("db connnection str:%s\n", connStr)

	dbClient, err := ent.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	db.Client = dbClient
	db.Ctx = context.Background()

	// auto migration
	if err := db.Client.Schema.Create(db.Ctx); err != nil {
		panic(err)
	}
	//Initializing redis

	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	db.Redis = redis.NewClient(&redis.Options{
		Addr: dsn, //redis port
	})
	_, err = db.Redis.Ping(db.Ctx).Result()
	if err != nil {
		panic(err)
	}
}
