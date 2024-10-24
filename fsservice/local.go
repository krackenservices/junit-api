package fsservice

import (
	"fmt"
	"log"
	"os"
)

type LocalFileSystem struct{}

func (lfs *LocalFileSystem) ListFiles(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
		return []string{}, err
	}
	filelist := []string{}
	for _, file := range files {
		fmt.Println(file.Name(), file.IsDir())
		if !file.IsDir() {
			filelist = append(filelist, file.Name())
		}
	}
	return filelist, nil
}

func (lfs *LocalFileSystem) GetFileContents(path string) ([]byte, error) {
	return os.ReadFile(path)
}
