package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"

	filepb "github.com/fuskovic/server-streaming/proto"

	"google.golang.org/grpc"
)

var filesDir = "files"

type server struct{}

func (s *server) Download(req *filepb.FileRequest, stream filepb.FileService_DownloadServer) error {
	fileName := req.GetFileName()
	path := filepath.Join(filesDir, fileName)

	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	fileSize := fileInfo.Size()

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var totalBytesStreamed int64

	for totalBytesStreamed < fileSize {
		shard := make([]byte, 1024)
		bytesRead, err := f.Read(shard)
		if err == io.EOF {
			log.Print("download complete")
			break
		}

		if err != nil {
			return err
		}

		if err := stream.Send(&filepb.FileResponse{
			Shard: shard,
		}); err != nil {
			return err
		}
		totalBytesStreamed += int64(bytesRead)
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on 50051 : %v\n", err)
	}

	s := grpc.NewServer()
	filepb.RegisterFileServiceServer(s, &server{})

	fmt.Println("starting gRPC server on 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to start server : %v\n", err)
	}
}
