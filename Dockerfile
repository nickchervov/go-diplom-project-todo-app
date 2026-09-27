FROM golang:1.26.5 AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /todo-build ./cmd/server/main.go

FROM alpine:latest

WORKDIR /app

RUN mkdir -p pkg/db

COPY --from=builder ./build/pkg/migrations ./pkg/migrations

COPY --from=builder ./todo-build ./todo-server

COPY --from=builder ./build/web ./web

ENV TODO_PORT=7540

ENV TODO_DBFILE=./pkg/db/scheduler.db

CMD [ "./todo-server" ]
