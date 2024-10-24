package fsservice

import "errors"

type FileSystemService struct {
	FsType string
	Name   string
	Fsi    FileSystemInterface
}

type FileSystemInterface interface {
	ListFiles(path string) ([]string, error)
	GetFileContents(path string) ([]byte, error)
}

func (f *FileSystemService) Init() error {
	switch f.FsType {
	case "local":
		f.Fsi = &LocalFileSystem{}
		return nil
	case "s3":
		s3fs, err := NewS3FileSystem(f.Name)
		if err != nil {
			return err
		}
		f.Fsi = s3fs
		return nil
	default:
		return errors.New("file system not supported")
	}
}
