build:
	GOARCH=amd64 GOOS=linux go build -o ./bin/app

deploy: build
	serverless deploy -v