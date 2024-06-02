package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database interface {
	getDB() *sql.DB
}

type database struct {
	db *sql.DB
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
