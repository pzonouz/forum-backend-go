run:
	@nodemon --exec "cd ./cmd/api && go run ." --ext "*.go" --signal SIGTERM
create-docker:
	docker run -d --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret  -e PGDATA=/var/lib/postgresql/data/pgdata -v ~/pgdata:/var/lib/postgresql/data postgres:16-alpine
remove-docker:
	docker stop postgres;docker rm postgres;
create-db:
	docker exec -it postgres createdb forum_go
drop-db:
	docker exec -it postgres dropdb forum_go
test:
	go test ./... 

testByCoverage:
	go test ./... -v -cover 

.PHONY: test clean all
