FROM golang:1.23.0-bookworm AS builder

WORKDIR /app

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server

FROM scratch AS prod
COPY --from=builder /app/server /server
COPY --from=builder /app/resources /resources
COPY --from=builder /app/migrations /app/migrations

ENTRYPOINT ["/server"]

FROM builder AS dev

RUN go install github.com/air-verse/air@latest
