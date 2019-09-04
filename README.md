# Server-Streaming gRPC API

Example code I wrote to demonstrate how to implement file downloads with a server-streaming gRPC API for my [blog](https://farishuskovic.dev/blog/server-streaming/).

### Dependencies and Pre-reqs

Install dependencies

    go get -u ./...

Create a test file in `server/files/`

    mkfile -v 100m server/files/someFileName

or just add some existing file you already have to `server/files/`

### Running

Open two shells one for the server and one for the client.

Run the server

`server/server.go`

    go run server.go

Run the client and pass it the filename of the file you added to `server/files/` in the pre-reqs step

`client/client.go`

    go run client.go someFileName

Example of successful client output:

        go run client.go testvideo.mp4
        23 MB downloaded   
        successfully downloaded testvideo.mp4

### Running the server with Docker

If your docker daemon is running I added a make file with some targets for running the file-server in a container and initiating a download.

Add any files you want to play with to `server/files/` before building the image with the first make target.

    make image

    make container

    make download file=testfile

The last make target is an example of a download using the testfile this repo comes with as the arg value. It downloads

To tear down the container use

    make stop