compose:
	docker compose -f docker-compose.yml up -d

test:
	go test -v ./...

run:
	go run cmd/tournament-service/main.go

binary:
	go build ./cmd/tournament-service

front:
	npm run build

