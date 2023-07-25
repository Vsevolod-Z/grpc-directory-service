package main

import (
	"context"
	"fmt"
	pb "grpc-directory-service/api/directory"
	c "grpc-directory-service/internal/cache"
	"grpc-directory-service/internal/server"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
)

type directoryServer struct {
	pb.UnsafeDirectoryServiceServer
}

var port = ":8089"

var cache *c.Cache

func (s *directoryServer) GetDirectoryInfo(ctx context.Context, req *pb.DirectoryRequest) (*pb.DirectoryResponse, error) {
	absolutePath := req.GetPath()
	_, err := os.Stat(absolutePath)
	if err != nil {
		if os.IsNotExist(err) {
			dirResponse := &pb.DirectoryResponse{
				Exists: false,
				Files:  []*pb.FileInfo{},
			}
			log.Printf("Path not exist: %v", absolutePath)
			return dirResponse, nil
		}
		return nil, err
	}
	if cacheData, ok := cache.Get(absolutePath); ok {

		dirResponse := &pb.DirectoryResponse{
			Exists: true,
			Files:  cacheData,
		}

		return dirResponse, nil
	}
	files, err := server.GetDirectoryInfo(absolutePath)
	if err != nil {
		return nil, err
	}
	dirResponse := &pb.DirectoryResponse{
		Exists: true,
		Files:  files,
	}
	cache.Add(absolutePath, files)
	return dirResponse, nil
}

func main() {
	server := grpc.NewServer()
	var directoryServer directoryServer
	pb.RegisterDirectoryServiceServer(server, &directoryServer)
	listen, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
		return
	}
	cache = c.NewCache(30*time.Second, 500)
	fmt.Println("Serving request...")
	server.Serve(listen)
}
