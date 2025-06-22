# # Build Stage 
# FROM golang:1.23-alpine as builder 
# WORKDIR /app
# COPY . .
# RUN go build -o main main.go
# RUN apk add --no-cache curl  

# # Download migrate
# RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.18.1/migrate.linux-amd64.tar.gz | tar xvz
 

# # Run Stage 
# FROM alpine 
# WORKDIR /app
# COPY --from=builder /app/main .
# COPY --from=builder /app/migrate ./migrate
# COPY app.env .
# COPY start.sh .
# # COPY wait-for.sh .
# COPY db/migration ./migration

# EXPOSE 8080
# ENTRYPOINT ["/app/start.sh"]
FROM golang:1.23 AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

FROM public.ecr.aws/lambda/provided:al2023
COPY --from=builder /app/main /main
COPY app.env /var/task/
ENTRYPOINT [ "/main" ]
