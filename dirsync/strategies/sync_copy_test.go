package strategies

import (
	"errors"
	"os"
	"testing"
)

func Test_syncFile(t *testing.T) {
	mockFileUtils := new(MockFileUtils)
	sc := &syncCopy{
		fileUtils: mockFileUtils,
	}
	cleanupMock := func() {
		mockFileUtils.ExpectedCalls = nil
		mockFileUtils.Calls = nil
	}

	source := "/source/file.txt"
	destination := "/destination/file.txt"

	t.Run("Destination file does not exist", func(t *testing.T) {
		mockFileUtils.On("Exists", destination).Return(false, nil)
		mockFileUtils.On("CopyFile", source, destination).Return(nil)

		sc.syncFile(source, destination)

		mockFileUtils.AssertCalled(t, "Exists", destination)
		mockFileUtils.AssertCalled(t, "CopyFile", source, destination)
		cleanupMock()
	})

	t.Run("Destination file exists and is the same", func(t *testing.T) {
		mockFileUtils.On("Exists", destination).Return(true, nil)
		mockFileUtils.On("IsTheSame", source, destination).Return(true, nil)

		sc.syncFile(source, destination)

		mockFileUtils.AssertCalled(t, "Exists", destination)
		mockFileUtils.AssertCalled(t, "IsTheSame", source, destination)
		mockFileUtils.AssertNotCalled(t, "CopyFile")
		cleanupMock()
	})

	t.Run("Destination file exists but is different", func(t *testing.T) {
		mockFileUtils.ExpectedCalls = nil
		mockFileUtils.On("Exists", destination).Return(true, nil)
		mockFileUtils.On("IsTheSame", source, destination).Return(false, nil)
		mockFileUtils.On("CopyFile", source, destination).Return(nil)

		sc.syncFile(source, destination)

		mockFileUtils.AssertCalled(t, "Exists", destination)
		mockFileUtils.AssertCalled(t, "IsTheSame", source, destination)
		mockFileUtils.AssertCalled(t, "CopyFile", source, destination)
		cleanupMock()
	})

	t.Run("Error checking if destination file exists", func(t *testing.T) {
		mockFileUtils.ExpectedCalls = nil
		mockFileUtils.Calls = nil
		mockFileUtils.On("Exists", destination).Return(false, errors.New("error checking existence"))

		sc.syncFile(source, destination)

		mockFileUtils.AssertCalled(t, "Exists", destination)
		mockFileUtils.AssertNotCalled(t, "IsTheSame", source, destination)
		mockFileUtils.AssertNotCalled(t, "CopyFile", source, destination)
		cleanupMock()
	})

	t.Run("Error checking if files are the same", func(t *testing.T) {
		mockFileUtils.ExpectedCalls = nil
		mockFileUtils.On("Exists", destination).Return(true, nil)
		mockFileUtils.On("IsTheSame", source, destination).Return(false, errors.New("error comparing files"))
		mockFileUtils.On("CopyFile", source, destination).Return(errors.New("error copying file"))

		sc.syncFile(source, destination)

		mockFileUtils.AssertCalled(t, "Exists", destination)
		mockFileUtils.AssertCalled(t, "IsTheSame", source, destination)
		mockFileUtils.AssertNotCalled(t, "CopyFile", source, destination)
		cleanupMock()
	})

	t.Run("Error copying the file", func(t *testing.T) {
		mockFileUtils.ExpectedCalls = nil
		mockFileUtils.On("Exists", destination).Return(false, nil)
		mockFileUtils.On("CopyFile", source, destination).Return(errors.New("error copying file"))

		sc.syncFile(source, destination)

		mockFileUtils.AssertCalled(t, "Exists", destination)
		mockFileUtils.AssertCalled(t, "CopyFile", source, destination)
		cleanupMock()
	})
}

func Test_syncDirectory(t *testing.T) {
	mockFileUtils := new(MockFileUtils)
	sc := &syncCopy{
		fileUtils: mockFileUtils,
	}
	cleanupMock := func() {
		mockFileUtils.ExpectedCalls = nil
		mockFileUtils.Calls = nil
	}
	source := "/source/dir"
	destination := "/destination/dir"

	t.Run("Destination directory exists, no recreating run", func(t *testing.T) {
		mockFileUtils.On("Exists", destination).Return(true, nil)
		mockFileUtils.On("RecreateDirectory", source, destination).Return(nil)
		mockFileUtils.On("ReadDirectory", source).Return([]os.DirEntry{}, errors.New("read directory error"))

		sc.syncDirectory(source, destination)

		mockFileUtils.AssertCalled(t, "Exists", destination)
		mockFileUtils.AssertNotCalled(t, "RecreateDirectory", source, destination)
		cleanupMock()
	})

	t.Run("Error Checking Destination", func(t *testing.T) {
		mockFileUtils.On("Exists", destination).Return(false, errors.New("error checking destination"))

		sc.syncDirectory(source, destination)

		mockFileUtils.AssertCalled(t, "Exists", destination)
		mockFileUtils.AssertNotCalled(t, "RecreateDirectory", source, destination)
		cleanupMock()
	})

	t.Run("Error Creating Destination", func(t *testing.T) {
		mockFileUtils.ExpectedCalls = nil
		mockFileUtils.On("Exists", destination).Return(false, nil)
		mockFileUtils.On("RecreateDirectory", source, destination).Return(errors.New("error creating directory"))

		sc.syncDirectory(source, destination)

		mockFileUtils.AssertCalled(t, "Exists", destination)
		mockFileUtils.AssertCalled(t, "RecreateDirectory", source, destination)
		cleanupMock()
	})
}
