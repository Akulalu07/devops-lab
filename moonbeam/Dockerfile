FROM golang:1.25.3-trixie
WORKDIR /moonbeam
COPY . .
RUN go mod tidy && go build -o moonbeam .

EXPOSE 8080
ENTRYPOINT ["./moonbeam"]
