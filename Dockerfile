# Build Stage 
FROM golang:1.22.0-alpine AS builder 
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add --no-cache curl  

# Download migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz
 

# Run Stage 
FROM alpine:latest 
WORKDIR /app

# 添加基础工具
RUN apk add --no-cache bash dos2unix

COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY db/migration ./migration

# 设置执行权限并转换行尾
RUN dos2unix /app/start.sh && \
    dos2unix /app/app.env && \
    chmod +x /app/start.sh

EXPOSE 8080
ENTRYPOINT ["/app/start.sh"]
