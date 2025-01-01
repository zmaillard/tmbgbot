FROM golang:1.23 as db-builder

# Install Goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Build Database
ADD db/ /db/

RUN GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=./tmbg.db goose -dir /db/sql/migrations up
