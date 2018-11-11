FROM golang:1.11.2

WORKDIR /transformer
ADD . /transformer
RUN go mod download; go build

ENTRYPOINT ["/transformer/transformer-go"]

# Expose the server TCP port
EXPOSE 8082