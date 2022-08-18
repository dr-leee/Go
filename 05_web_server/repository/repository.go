package repository

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type DbModel struct {
	DB *pgxpool.Pool
}

type Repo interface {
	GetLogs(ctx context.Context) (pgx.Rows, error)
	CreateLog(ctx context.Context, login string) error
}

func (conn *DbModel) GetLogs(ctx context.Context) (pgx.Rows, error) {

	rows, err := conn.DB.Query(ctx, "SELECT * FROM logins")
	if err != nil {
		log.Print("Unable to FETCH: %v\n", err)
	}
	return rows, err
}

func (conn *DbModel) CreateLog(ctx context.Context, login string) error {
	//err := conn.DB.Ping(conn.CTX)

	_, err := conn.DB.Query(ctx, "INSERT INTO logins(login, time) VALUES ($1, now())", login)
	if err != nil {
		log.Print("Unable to INSERT: %v\n", err)
	}
	return err
}
