# Makefile

# 1. The equivalent of `npm run dev`
dev:
	air

# 2. The equivalent of `npm run build`
build:
	go build -o bin/api cmd/main.go

# 3. Start the production binary
start:
	./bin/api

# 4. Generate sqlc code
db-gen:
	sqlc generate

# 5. Run tests
test:
	go test -v ./...
