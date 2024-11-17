compose:
	docker compose -f docker-compose.yml up -d

run:
	go run cmd/tournament-service/main.go

binary:
	go build ./cmd/tournament-service

front:
	npm run build

