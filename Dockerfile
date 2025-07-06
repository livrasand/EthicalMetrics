FROM golang:1.23-bullseye

RUN apt-get update && apt-get install -y pkg-config && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .

RUN go mod tidy

RUN go build -tags "sqlcipher" -o server ./cmd/server

EXPOSE 8080

CMD ["./server"]
