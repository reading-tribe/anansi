
   
gitHash = $(shell git rev-parse HEAD)

build:
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_create_book func_create_book/*.go
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_delete_book func_delete_book/*.go
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_get_book func_get_book/*.go
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_list_books func_list_books/*.go
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_update_book func_update_book/*.go
	
.PHONY: clean
clean:
	rm -rf ./bin
	mkdir bin

.PHONY: deploy
deploy: clean build
	./deploy.sh