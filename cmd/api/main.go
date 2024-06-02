package main

func main() {
	router := NewRouter()
	server := NewServer(router.mux)
	server.serve()
}
