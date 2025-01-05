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
COPY tmbg.db /tmbg.db

USER nonroot:nonroot

ENV GOOSE_DRIVER=sqlite3
ENV GOOSE_DBSTRING=./tmbg.db
ENTRYPOINT ["/tmbgbot"]

