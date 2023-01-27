BINARY_NAME = app
.DEFAULT_GOAL = run

build:
	GOARCH=amd64 GOOS=linux   go build -o ./bin/${BINARY_NAME}-linux       cmd/main/app.go
	GOARCH=amd64 GOOS=windows go build -o ./bin/${BINARY_NAME}-windows.exe cmd/main/app.go

run:
	go run ./cmd/main/app.go

clean:
	go clean
	rm ./bin/*