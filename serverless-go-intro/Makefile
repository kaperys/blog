.PHONY: build

fmt:
	go fmt 
	go vet

test: fmt
	go test -cover -race

build:
	GOOS=linux go build -o handler

sam: build
	sam local start-api
