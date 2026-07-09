FROM golang:1.25.11-alpine AS builder

WORKDIR /app 
COPY go.mod go.sum ./ 
RUN go mod download 
COPY . . 
RUN go build -o /vatify ./cmd/api/ 
FROM alpine
WORKDIR /app 
COPY --from=builder /vatify /vatify
COPY internal/html/index.html /app/index.html
COPY internal/html/success.html /app/success.html
COPY internal/html/cancel.html /app/cancel.html
COPY img/icon.png /app/icon.png
COPY img/s1.png /app/s1.png
COPY img/s2.png /app/s2.png
COPY img/s3.png /app/s3.png
COPY img/s4.png /app/s4.png
EXPOSE 8080
CMD ["/vatify"]