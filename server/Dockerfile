FROM golang:1.23-alpine AS builder

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /src


# 复制全部源码（在构建上下文为仓库根时）
COPY . .
# 下载依赖（利用 go.mod 缓存）
RUN go mod tidy

# 构建二进制
RUN go build ./...

ENTRYPOINT ["go", "run", "."]
