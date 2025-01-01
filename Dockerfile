FROM golang:1.23 AS db-builder

# Install Goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /app

# Build Database
COPY db/ ./db

RUN GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=./tmbg.db goose -dir ./db/sql/migrations up

FROM golang:1.23 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY dbstore/ ./dbstore

RUN CGO_ENABLED=0 GOOS=linux go build -o /tmbgbot

FROM gcr.io/distroless/base-debian11 AS app

WORKDIR /
COPY --from=builder /tmbgbot /tmbgbot
COPY --from=db-builder /app/tmbg.db /tmbg.db

USER nonroot:nonroot

ENV GOOSE_DRIVER=sqlite3
ENV GOOSE_DBSTRING=./tmbg.db
ENTRYPOINT ["/tmbgbot"]

