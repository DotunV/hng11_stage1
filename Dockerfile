FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o ./example-golang ./main.go
 
 
FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/example-golang .
EXPOSE 8080
ENTRYPOINT ["./"]