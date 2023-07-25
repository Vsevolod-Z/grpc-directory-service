package server

import (
	"fmt"
	pb "grpc-directory-service/api/directory"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func GetDirectoryInfo(absolutePath string) ([]*pb.FileInfo, error) {
	files, err := ioutil.ReadDir(absolutePath)
	if err != nil {
		log.Fatal(err)
	}
	filesInfo := make([]*pb.FileInfo, 0, 10)

	for _, file := range files {

		fileInfo := &pb.FileInfo{}
		if file.IsDir() {
			fileInfo.Name = file.Name()
			fileInfo.Size, err = getDirectorySize(fmt.Sprintf("%s/%s", absolutePath, file.Name()))
			if err != nil {
				log.Println(err)
			}

		} else {
			fileInfo.Name = file.Name()
			fileInfo.Size = file.Size()
		}
		fmt.Println(fileInfo)
		filesInfo = append(filesInfo, fileInfo)
	}
	return filesInfo, err

}

func getDirectorySize(path string) (int64, error) {
	var size int64
	err := filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fileInfo.IsDir() {
			size += fileInfo.Size()
		}
		return nil
	})
	return size, err
}
