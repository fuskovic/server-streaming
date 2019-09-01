package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	//using humanize to format B,MB,GB, etc...
	"github.com/dustin/go-humanize"
	filepb "github.com/fuskovic/server-streaming/proto"

	"google.golang.org/grpc"
)

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatalf("Please provide a filename argument")
	}

	requestedFile := args[0]

	cc, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to establish connection with gRPC server : %v\n", err)
	}
	defer cc.Close()

	c := filepb.NewFileServiceClient(cc)

	if err := download(requestedFile, c); err != nil {
		log.Fatalf("failed to download %s : %v\n", requestedFile, err)
	}
	fmt.Printf("\nsuccessfully downloaded %s\n", requestedFile)
}

func download(fileName string, client filepb.FileServiceClient) error {
	req := &filepb.FileRequest{
		FileName: fileName,
	}

	stream, err := client.Download(context.Background(), req)
	if err != nil {
		return err
	}

	var downloaded int64
	var buffer bytes.Buffer

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			if err := ioutil.WriteFile(fileName, buffer.Bytes(), 0777); err != nil {
				return err
			}
			break
		}
		if err != nil {
			buffer.Reset()
			return err
		}

		shard := res.GetShard()
		shardSize := len(shard)
		downloaded += int64(shardSize)

		buffer.Write(shard)
		fmt.Printf("\r%s", strings.Repeat(" ", 25))
		fmt.Printf("\r%s downloaded", humanize.Bytes(uint64(downloaded)))
	}
	return nil
}
