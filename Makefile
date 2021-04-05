all: test install lint

install:
	go install ./...

test:
	go test ./...

lint:
	golangci-lint run

artifacts: artifacts-linux-x64 artifacts-windows-x64 artifacts-arm-64

artifacts-linux-x64:
	@echo "-- build linux x64"
	@GOOS=linux GOARCH=amd64 go build -o artifacts-linux-x64/ ./...

artifacts-windows-x64:
	@echo "-- build windows x64"
	@GOOS=windows GOARCH=amd64 go build -o artifacts-windows-x64/ ./...

artifacts-arm-64:
	@echo "-- build arm64"
	@GOOS=linux GOARCH=arm64 go build -o artifacts-arm-x64/ ./...
