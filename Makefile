all: test install lint

install:
	go install ./...

test:
	go test ./...

lint:
	golangci-lint run
