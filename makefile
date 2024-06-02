run:
	@nodemon --exec "cd ./cmd/api && go run ." --ext "*.go" --signal SIGTERM
test:
	go test ./... 
