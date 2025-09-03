FROM golang:1.24.6-bookworm

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/bot
RUN CGO_ENABLED=0 GOOS=linux go build -o bot

WORKDIR /app
CMD ["./cmd/bot/bot"]
