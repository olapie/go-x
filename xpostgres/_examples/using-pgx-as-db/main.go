package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.olapie.com/x/xpostgres"
)

func main() {
	pool, err := pgxpool.New(context.TODO(), xpostgres.NewOpenOptions().String())
	if err != nil {
		log.Fatalln(err)
	}
	db := stdlib.OpenDBFromPool(pool)
	err = db.PingContext(context.TODO())
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println("connected")
	}
}
