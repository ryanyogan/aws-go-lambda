default: build

test:
	go get -d -t ./... && go vet -x ./... && go test ./...

build:
	dep ensure
	rm -rf bin/
	env GOOS=linux go build -ldflags="-s -w" -o bin/handlers/addTodo src/handlers/create/addTodo.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/handlers/listTodos src/handlers/list/listTodos.go

release: test build
	serverless deploy