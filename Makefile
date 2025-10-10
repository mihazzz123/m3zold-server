BINARY_NAME=m3zold-server
CMD_PATH=./cmd

build:
	go build -v -o $(BINARY_NAME) $(CMD_PATH)

build-debug:
	go build -gcflags="all=-N -l" -v -o $(BINARY_NAME) $(CMD_PATH)

run:
	./$(BINARY_NAME)

test:
	go test ./...

up:
	docker-compose up --build

down:
	docker-compose down

debug-up:
	docker-compose -f docker-compose.debug.yml up --build

debug-down:
	docker-compose -f docker-compose.debug.yml down

clean:
	docker-compose down --volumes --remove-orphans
	docker system prune -af
	rm -f $(BINARY_NAME)
