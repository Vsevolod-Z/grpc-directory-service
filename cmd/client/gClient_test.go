package main

import (
	"context"
	pb "grpc-directory-service/api/directory"
	"testing"

	"google.golang.org/grpc"
)

func TestGetDirectoryInfo(t *testing.T) {
	conn, err := grpc.Dial("localhost:8089", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	client := pb.NewDirectoryServiceClient(conn)

	ctx := context.Background()
	path := "S:\\ArkServerManager"
	response, err := client.GetDirectoryInfo(ctx, &pb.DirectoryRequest{Path: path})
	if err != nil {
		t.Fatalf("Failed to call GetDirectoryInfo: %v", err)
	}

	if !response.Exists {
		t.Errorf("Path not exist: %v", path)
	}
}
