package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"

	"forum-backend-go/ineternal/utils"
)

type Database interface {
	getDB() *sql.DB
	InitDB()
	RunQueryOnDB(string)
}

type database struct {
	db *sql.DB
}

func (d *database) InitDB() {
	d.RunQueryOnDB(utils.CreateUserTableQuery)

}

func (d *database) RunQueryOnDB(query string) {
	_, err := d.db.Exec(query)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func NewDatabase() *database {
	dsn := GetEnv("dsn", "host=localhost port=5432 user=root password=secret dbname=forum sslmode=disable")
	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err.Error())
	}
	err = conn.Ping()
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Succesfully connect to POSTGRESQL!")
	return &database{db: conn}
}

func (d *database) getDB() *sql.DB {
	return d.db
}
