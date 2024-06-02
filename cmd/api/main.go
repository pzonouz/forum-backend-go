package main

func main() {

	var database Database = NewDatabase()
	db := database.getDB()
	defer db.Close()
	var router Router = NewRouter(db)
	router.registerRoutes(db)
	NewServer(router.getMux()).serve()
}
