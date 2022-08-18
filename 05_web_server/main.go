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

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	env := &rest.Env{
		Repo: repository.DbModel{DB: db},
	}

	http.HandleFunc("/login", env.GetJsonLogin)
	http.HandleFunc("/showlogins", env.ShowLogins)
	log.Println(http.ListenAndServe(":8090", nil))
}
