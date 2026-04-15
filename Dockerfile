# ── 阶段一：编译 ─────────────────────────────────────────────────
FROM golang:1.20-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s" -o ./app main.go wire_gen.go

# ── 阶段二：运行 ─────────────────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache tzdata ca-certificates \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/application.prod.yml .

RUN mkdir -p logs

EXPOSE 8000

ENTRYPOINT ["./app", "-envString", "prod"]
