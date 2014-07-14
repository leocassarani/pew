NAME=pew

build:
	@mkdir -p bin/
	go build -o bin/$(NAME)

fmt:
	go fmt ./...

test:
	go test ./...
