BINARY := game
CMD_PATH := ./cmd/tankismus

.PHONY: build run test clean

build:
	go build -o bin/$(BINARY) $(CMD_PATH)

run: build
	./bin/$(BINARY)

test:
	go test ./...

clean:
	rm -f bin/$(BINARY)
