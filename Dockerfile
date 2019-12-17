FROM golang:latest

LABEL maintainer="cydu@google.com>"
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server .

EXPOSE 9102
CMD ["./server"]
