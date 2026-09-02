FROM golang:1.26.5
WORKDIR /app
COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg
COPY go.sum .

RUN apt update -y && apt install ffmpeg -y
RUN go build -o bot cmd/main.go
CMD ["./bot"]