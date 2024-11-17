FROM golang:1.23-alpine AS golangbuild

WORKDIR /app

COPY ../ ./

RUN go mod tidy && go mod vendor
RUN go build ./cmd/tournament-service

FROM alpine:3.14

WORKDIR /app

COPY --from=golangbuild /app/tournament-service ./tournament-service
COPY --from=golangbuild /app/assets ./assets
COPY --from=golangbuild /app/templates ./templates

CMD ["./tournament-service"]