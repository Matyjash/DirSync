package fileutils

import "os"

type FileUtils interface {
	Exists(path string) (bool, error)
	CopyFile(sourceFilePath, destinationFilePath string) error
	RecreateDirectory(sourceDir, destinationDir string) error
	IsTheSame(sourceFilePath, destinationFilePath string) (bool, error)
	ReadDirectory(path string) ([]os.DirEntry, error)
	DeleteAll(path string) error
}
