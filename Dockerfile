# Build Stage 
FROM golang:1.23-alpine as builder 
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run Stage 
FROM alpine 
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .

EXPOSE 8080
CMD ["/app/main"]