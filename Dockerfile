FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
ENV GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . .
RUN go build -o /worker ./cmd/worker

FROM debian:bookworm-slim
RUN apt-get update \                                                                                                                                                       
      && apt-get install -y --no-install-recommends ca-certificates \                                                                                                        
      && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /worker /usr/local/bin/worker

ENTRYPOINT ["/usr/local/bin/worker"]
