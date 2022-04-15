
   
gitHash = $(shell git rev-parse HEAD)

build:
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_register func_register/*.go
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_login func_login/*.go
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_accept_invite func_accept_invite/*.go
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_invite_user func_invite_user/*.go
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_logout func_logout/*.go
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_refresh func_refresh/*.go
	env GOARCH=amd64 GOOS=linux go build -trimpath -ldflags="-s -w" -o bin/func_authorizer func_authorizer/*.go
	
.PHONY: clean
clean:
	rm -rf ./bin
	mkdir bin

.PHONY: deploy
deploy: clean build
	./deploy.sh