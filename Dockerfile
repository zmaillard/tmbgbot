ARG XCPUTRANSLATE_VERSION=v0.6.0
ARG BUILDPLATFORM=linux/amd64
FROM --platform=${BUILDPLATFORM} qmcgaw/xcputranslate:${XCPUTRANSLATE_VERSION} AS xcputranslate
FROM --platform=${BUILDPLATFORM}  golang:1.23 AS builder

WORKDIR /app
ARG TARGETPLATFORM

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY dbstore/ ./dbstore

ENV CGO_ENABLED=0

RUN GOARCH="$(xcputranslate translate -field arch -targetplatform ${TARGETPLATFORM})" \
    GOARM="$(xcputranslate translate -field arm -targetplatform ${TARGETPLATFORM})" \
    go build -o /tmbgbot

FROM gcr.io/distroless/base-debian11 AS app
ARG TARGETPLATFORM
WORKDIR /
COPY --from=builder /tmbgbot /tmbgbot
COPY tmbg.db /tmbg.db

USER nonroot:nonroot

ENV GOOSE_DRIVER=sqlite3
ENV GOOSE_DBSTRING=./tmbg.db
ENTRYPOINT ["/tmbgbot"]

