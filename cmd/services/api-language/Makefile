
   
gitHash = $(shell git rev-parse HEAD)

build:
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_list_languages func_list_languages/*.go

.PHONY: clean
clean:
	rm -rf ./bin
	mkdir bin

.PHONY: deploy
deploy: clean build
	./deploy.sh