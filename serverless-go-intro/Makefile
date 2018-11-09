S3_BUCKET = serverless-go-intro
STACK_NAME = ServelessGoIntroStack

.PHONY: build

fmt:
	go fmt && go vet

test: fmt
	go test -cover -race

build:
	GOOS=linux GOARCH=amd64 go build -o handler

validate:
	sam validate

api: build
	sam local start-api

invoke: build
	sam local invoke "handler" -e invoke.json

package: build
	sam package --template-file template.yml --s3-bucket $(S3_BUCKET) --output-template-file package.yml

deploy: package
	sam deploy --stack-name $(STACK_NAME) --template-file package.yml --capabilities CAPABILITY_IAM
