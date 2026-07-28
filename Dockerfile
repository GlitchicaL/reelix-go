# Dockerfile

FROM golang:1.26-alpine

RUN apk add --no-cache ffmpeg

WORKDIR /app

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o reelix-go ./cmd/reelix-go

CMD ["./reelix-go"]