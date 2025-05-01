FROM golang:1.24.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o TestCloud- ./cmd
FROM gcr.io/distroless/base-debian11

WORKDIR /root/
COPY --from=builder /app/TestCloud- .
COPY --from=builder /app/migration ./migration
COPY config ./config
EXPOSE 8081

CMD ["./TestCloud-","-config", "config/config.yaml"]