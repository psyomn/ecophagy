all: test install lint

install:
	go install ./...

test:
	go test ./...

lint:
	golangci-lint run

artifacts: artifacts-linux-x64 artifacts-windows-x64 # artifacts-arm-64

artifacts-linux-x64:
	@echo "-- build linux x64"
	@GOOS=linux GOARCH=amd64 go build -o artifacts-linux-x64/ ./...

# you need multilib support installed for this
artifacts-windows-x64:
	@echo "-- build windows x64"
	@GOOS=windows GOARCH=386 CGO_ENABLED=1 CXX=i686-w64-mingw32-g++ CC=i686-w64-mingw32-gcc go build -o artifacts-windows-x64/ ./...

# TODO: need to study how to enable this
# artifacts-arm-64:
# 	@echo "-- build arm64"
# 	@GOOS=linux GOARCH=arm64 CGO_ENABLED=1 go build -o artifacts-arm-x64/ ./...
