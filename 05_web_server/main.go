package main

import (
	"context"
	"github.com/dr-leee/Go/05_web_server/delivery/rest"
	"github.com/dr-leee/Go/05_web_server/repository"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
)

const dbString = "postgres://test:test@localhost:6000/postgres?sslmode=disable"

func main() {

	db, err := pgxpool.Connect(context.Background(), dbString)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = db.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	env := &rest.Env{
		Repo: repository.DbModel{DB: db, CTX: ctx},
	}

	http.HandleFunc("/login", env.GetJsonLogin)
	log.Println(http.ListenAndServe(":8090", nil))
}
