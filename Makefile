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
	docker-compose -f docker-compose.debug.yml up

debug-restart:
	docker-compose -f docker-compose.debug.yml restart api-debug

debug-logs:
	docker-compose -f docker-compose.debug.yml logs -f api-debug

debug-down:
	docker-compose -f docker-compose.debug.yml down

# Полный перезапуск
debug-rebuild: debug-down
	docker-compose -f docker-compose.debug.yml up --build