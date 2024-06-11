run:
	@nodemon --exec "cd ./cmd/api && go run ." --ext "*.go" --signal SIGTERM
test:
	go test ./... 

testByCoverage:
	go test ./... -v -cover 

.PHONY: test clean all
