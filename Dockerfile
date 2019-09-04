FROM golang:latest

ENV dir src/github.com/fuskovic/server-streaming/
RUN mkdir -p $dir
WORKDIR $dir
COPY . .
RUN go get -u github.com/dustin/go-humanize \
github.com/golang/protobuf/proto \
google.golang.org/grpc
RUN go build -o server/run server/server.go