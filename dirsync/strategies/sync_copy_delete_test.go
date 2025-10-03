package strategies

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func Test_deleteNotInSource(t *testing.T) {
	mockFileUtils := new(MockFileUtils)
	scd := &SyncCopyDelete{
		fileUtils: mockFileUtils,
	}
	cleanupMock := func() {
		mockFileUtils.ExpectedCalls = nil
		mockFileUtils.Calls = nil
	}

	source := "/source/dir"
	destination := "/destination/dir"

	t.Run("Error reading destination directory", func(t *testing.T) {
		mockFileUtils.On("ReadDirectory", destination).Return([]fs.DirEntry{}, errors.New("error reading directory"))

		scd.deleteNotInSource(source, destination)

		mockFileUtils.AssertCalled(t, "ReadDirectory", destination)
		mockFileUtils.AssertNotCalled(t, "Exists")
		cleanupMock()
	})

	t.Run("File in destination but not in source", func(t *testing.T) {
		files := []os.DirEntry{}
		mockFile := new(MockDirEntry)
		mockFile.On("Name").Return("file1.txt")
		mockFile.On("IsDir").Return(false)
		files = append(files, mockFile)

		mockFileUtils.On("ReadDirectory", destination).Return(files, nil)
		mockFileUtils.On("Exists", filepath.Join(source, "file1.txt")).Return(false, nil)
		mockFileUtils.On("DeleteAll", filepath.Join(destination, "file1.txt")).Return(nil)

		scd.deleteNotInSource(source, destination)

		mockFileUtils.AssertCalled(t, "ReadDirectory", destination)
		mockFileUtils.AssertCalled(t, "Exists", filepath.Join(source, "file1.txt"))
		mockFileUtils.AssertCalled(t, "DeleteAll", filepath.Join(destination, "file1.txt"))
		cleanupMock()
	})

	t.Run("File in destination and in source", func(t *testing.T) {
		mockFile := new(MockDirEntry)
		mockFile.On("Name").Return("file1.txt")
		mockFile.On("IsDir").Return(false)

		mockFileUtils.On("ReadDirectory", destination).Return([]os.DirEntry{mockFile}, nil)
		mockFileUtils.On("Exists", filepath.Join(source, "file1.txt")).Return(true, nil)

		scd.deleteNotInSource(source, destination)

		mockFileUtils.AssertCalled(t, "ReadDirectory", destination)
		mockFileUtils.AssertCalled(t, "Exists", filepath.Join(source, "file1.txt"))
		mockFileUtils.AssertNotCalled(t, "DeleteAll", filepath.Join(destination, "file1.txt"))
		cleanupMock()
	})

	t.Run("Directory in destination and in source", func(t *testing.T) {
		mockDir := new(MockDirEntry)
		mockDir.On("Name").Return("subdir")
		mockDir.On("IsDir").Return(true)

		mockFileUtils.On("ReadDirectory", destination).Return([]os.DirEntry{mockDir}, nil).Once()
		mockFileUtils.On("ReadDirectory", filepath.Join(destination, "subdir")).Return([]os.DirEntry{}, nil).Once()
		mockFileUtils.On("Exists", filepath.Join(source, "subdir")).Return(true, nil)

		scd.deleteNotInSource(source, destination)

		mockFileUtils.AssertCalled(t, "ReadDirectory", destination)
		mockFileUtils.AssertCalled(t, "Exists", filepath.Join(source, "subdir"))
		mockFileUtils.AssertNotCalled(t, "DeleteAll", filepath.Join(destination, "subdir"))
		cleanupMock()
	})

	t.Run("Error checking if file/directory exists in source", func(t *testing.T) {
		mockFile := new(MockDirEntry)
		mockFile.On("Name").Return("file1.txt")
		mockFile.On("IsDir").Return(false)

		mockFileUtils.On("ReadDirectory", destination).Return([]os.DirEntry{mockFile}, nil)
		mockFileUtils.On("Exists", filepath.Join(source, "file1.txt")).Return(false, errors.New("error checking existence"))

		scd.deleteNotInSource(source, destination)

		mockFileUtils.AssertCalled(t, "ReadDirectory", destination)
		mockFileUtils.AssertCalled(t, "Exists", filepath.Join(source, "file1.txt"))
		mockFileUtils.AssertNotCalled(t, "DeleteAll", filepath.Join(destination, "file1.txt"))
		cleanupMock()
	})

	t.Run("Error deleting a file/directory", func(t *testing.T) {
		mockFile := new(MockDirEntry)
		mockFile.On("Name").Return("file1.txt")
		mockFile.On("IsDir").Return(false)

		mockFileUtils.On("ReadDirectory", destination).Return([]os.DirEntry{mockFile}, nil)
		mockFileUtils.On("Exists", filepath.Join(source, "file1.txt")).Return(false, nil)
		mockFileUtils.On("DeleteAll", filepath.Join(destination, "file1.txt")).Return(errors.New("error deleting file"))

		scd.deleteNotInSource(source, destination)

		mockFileUtils.AssertCalled(t, "ReadDirectory", destination)
		mockFileUtils.AssertCalled(t, "Exists", filepath.Join(source, "file1.txt"))
		mockFileUtils.AssertCalled(t, "DeleteAll", filepath.Join(destination, "file1.txt"))
	})
}
