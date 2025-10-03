package fileutils

import (
	"fmt"
	"io"
	"os"
)

type DefaultFileHandler struct{}

func NewDefaultFileHandler() *DefaultFileHandler {
	return &DefaultFileHandler{}
}

func (dfh *DefaultFileHandler) RecreateDirectory(sourceDir, destinationDir string) error {
	sourceInfo, err := os.Stat(sourceDir)
	if err != nil {
		return fmt.Errorf("failed to get source directory info: %w", err)
	}

	err = os.Mkdir(destinationDir, sourceInfo.Mode())
	if err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	return nil
}

func (dfh *DefaultFileHandler) IsTheSame(sourceFilePath, destinationFilePath string) (bool, error) {
	sourceInfo, err := os.Stat(sourceFilePath)
	if err != nil {
		return false, fmt.Errorf("failed to get source file info: %w", err)
	}

	destInfo, err := os.Stat(destinationFilePath)
	if err != nil {
		return false, fmt.Errorf("failed to get destination file info: %w", err)
	}

	if sourceInfo.Size() != destInfo.Size() {
		return false, nil
	}

	if !sourceInfo.ModTime().Equal(destInfo.ModTime()) {
		return false, nil
	}

	return true, nil
}

func (dfh *DefaultFileHandler) Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func (dfh *DefaultFileHandler) CopyFile(sourceFilePath, destinationFilePath string) error {
	sourceFile, err := os.Open(sourceFilePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destinationFilePath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	err = dfh.preserveModificationAndAccessTime(sourceFilePath, destinationFilePath)
	if err != nil {
		return fmt.Errorf("failed to preserve file times: %w", err)
	}
	return nil
}

func (dfh *DefaultFileHandler) ReadDirectory(path string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	return entries, nil
}

func (dfh *DefaultFileHandler) DeleteAll(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return fmt.Errorf("failed to delete path: %w", err)
	}
	return nil
}

func (dfh *DefaultFileHandler) preserveModificationAndAccessTime(sourceFilePath, destinationFilePath string) error {
	//Without this, copied files have current time as mod/access time and would be treated as different next time
	sourceInfo, err := os.Stat(sourceFilePath)
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	err = os.Chtimes(destinationFilePath, sourceInfo.ModTime(), sourceInfo.ModTime())
	if err != nil {
		return fmt.Errorf("failed to set file times: %w", err)
	}
	return nil
}
