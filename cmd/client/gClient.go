package main

import (
	"bufio"
	"context"
	"fmt"
	pb "grpc-directory-service/api/directory"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
)

var port = ":8089"

func main() {
	conn, err := grpc.Dial(fmt.Sprintf("%s%s", "localhost", port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewDirectoryServiceClient(conn)
	ctx := context.Background()
	for {
		path := readUserInput("Enter the directory path (or type 'exit' to quit): ")
		if path == "exit" {
			break
		}
		fmt.Println(path)
		response, err := client.GetDirectoryInfo(ctx, &pb.DirectoryRequest{Path: path})
		if err != nil {
			log.Printf("failed to call GetDirectoryInfo: %v", err)
		}
		if response.Exists {
			for _, fileInfo := range response.Files {
				log.Printf("File name: %s, Size: %.2f MB", fileInfo.Name, float64(fileInfo.Size)/(1024*1024))
			}
		} else {
			log.Printf("Path not exist: %v", path)
		}
	}
}
func readUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("failed to read user input: %v", err)
	}
	return strings.TrimSpace(input)
}
