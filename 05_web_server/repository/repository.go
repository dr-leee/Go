package repository

import (
	"context"
	"github.com/dr-leee/Go/05_web_server/types"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type DbModel struct {
	DB  *pgxpool.Pool
	CTX context.Context
}

type Repo interface {
	GetLogs() (types.Log, error)
	CreateLog(login string) error
}

/*
func (conn *DbModel) GetLogs() (types.Log, error) {

}
*/
func (conn *DbModel) CreateLog(login string) error {
	//err := conn.DB.Ping(conn.CTX)

	_, err := conn.DB.Query(conn.CTX, "INSERT INTO logins(login, time) VALUES ($1, now())", login)
	if err != nil {
		log.Print("Unable to INSERT: %v\n", err)
	}
	return err
}
