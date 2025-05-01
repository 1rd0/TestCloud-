# ---------- build ----------
FROM golang:1.24.2-alpine AS builder
WORKDIR /src

COPY go.mod .
RUN go mod download
COPY . .

# собираем все 4 бинарника
RUN for n in 1 2 3 4; do \
        CGO_ENABLED=0 GOOS=linux go build -o /out/testServer${n} testServer${n}.go ; \
    done

# ---------- runtime ----------
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /out/* /app/

# запускаем нужный файл по переменной BIN
RUN printf '#!/bin/sh\nexec /app/$BIN\n' > /start && chmod +x /start
ENV BIN=testServer1
EXPOSE 9001 9002 9003 9004
ENTRYPOINT ["/start"]
