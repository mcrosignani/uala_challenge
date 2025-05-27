ARG WORKDIR_DEFAULT=/go/src/github.com/mcrosignani/uala_challenge/users

FROM golang:latest

ARG WORKDIR_DEFAULT
WORKDIR $WORKDIR_DEFAULT

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/httpserver/main

ENV POSTGRES_MIGRATIONS_PATH="$WORKDIR_DEFAULT/db/migrations"

EXPOSE 8080

CMD ["./main"]