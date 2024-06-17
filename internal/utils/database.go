package utils

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database interface {
	GetDB(isTest bool) (*sql.DB, error)
	RunQueryOnDB(query string) error
	// QueryRow(tableName string)
}

type database struct {
	db *sql.DB
}

func (d *database) RunQueryOnDB(query string) error {
	_, err := d.db.Exec(query)

	return err
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

	log.Println("Successfully connect to POSTGRESQL!")

	return &database{db: conn}
}

func (d *database) GetDB(isTest bool) (*sql.DB, error) {
	if isTest {
		err := d.RunQueryOnDB(CreateRoleTableQueryTest)
		if err != nil {
			return nil, err
		}

		err = d.RunQueryOnDB(CreateUserTableQueryTest)

		if err != nil {
			return nil, err
		}

		return d.db, nil
	}

	err := d.RunQueryOnDB(CreateRoleTableQuery)

	if err != nil {
		return nil, err
	}

	err = d.RunQueryOnDB(CreateUserTableQuery)

	if err != nil {
		return nil, err
	}

	return d.db, nil
}

func (d *database) TearDown(tableName string) {
	err := d.RunQueryOnDB(fmt.Sprintf(DeleteTestTableQuery, tableName))
	if err != nil {
		log.Print(err.Error())
	}
}
