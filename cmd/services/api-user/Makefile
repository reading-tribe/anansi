
   
gitHash = $(shell git rev-parse HEAD)

build:
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_create_user func_create_user/*.go
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_delete_user func_delete_user/*.go
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_get_user func_get_user/*.go
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_list_users func_list_users/*.go
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_update_user func_update_user/*.go
	
.PHONY: clean
clean:
	rm -rf ./bin
	mkdir bin

.PHONY: deploy
deploy: clean build
	./deploy.sh