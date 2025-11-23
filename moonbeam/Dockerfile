FROM golang:1.25.4-trixie AS builder
WORKDIR /moonbeam
ENV GOEXPERIMENT=greenteagc
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -o moonbeam .

FROM debian:bookworm-slim
WORKDIR /moonbeam
RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
COPY --from=builder /moonbeam/moonbeam .
EXPOSE 8080
ENTRYPOINT ["./moonbeam"]
