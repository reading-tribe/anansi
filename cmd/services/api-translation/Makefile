
   
gitHash = $(shell git rev-parse HEAD)

build:
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_create_translation func_create_translation/*.go
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_delete_translation func_delete_translation/*.go
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_get_translation func_get_translation/*.go
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_list_translations func_list_translations/*.go
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_update_translation func_update_translation/*.go
	
.PHONY: clean
clean:
	rm -rf ./bin
	mkdir bin

.PHONY: deploy
deploy: clean build
	./deploy.sh