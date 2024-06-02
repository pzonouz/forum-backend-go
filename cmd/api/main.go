package main

func main() {
	var database Database = NewDatabase()

	database.InitDB()

	db := database.getDB()
	defer db.Close()

	var router Router = newRouter(db)

	router.registerRoutes(db)
	newServer(router.getMux()).serve()
}
