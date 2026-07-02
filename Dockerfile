FROM golang:1.25.11-alpine AS builder

WORKDIR /app 
COPY go.mod go.sum ./ 
RUN go mod download 
COPY . . 
RUN go build -o /vatify ./cmd/api/ 
FROM alpine
WORKDIR /app 
COPY --from=builder /vatify /vatify 
EXPOSE 8080
CMD ["/vatify"]