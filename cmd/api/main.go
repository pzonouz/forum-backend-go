package main

import (
	"log"

	"forum-backend-go/internal/utils"
)

func main() {
	var database utils.Database = utils.NewDatabase()

	db, err := database.GetDB(false)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	var router Router = newRouter(db)

	router.registerRoutes()
	newServer(router.getMux()).serve()
}
