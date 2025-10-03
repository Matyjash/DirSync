package strategies

import (
	"os"

	"github.com/stretchr/testify/mock"
)

type MockFileUtils struct {
	mock.Mock
}

func (m *MockFileUtils) Exists(path string) (bool, error) {
	args := m.Called(path)
	return args.Bool(0), args.Error(1)
}

func (m *MockFileUtils) RecreateDirectory(source, destination string) error {
	args := m.Called(source, destination)
	return args.Error(0)
}

func (m *MockFileUtils) CopyFile(source, destination string) error {
	args := m.Called(source, destination)
	return args.Error(0)
}

func (m *MockFileUtils) IsTheSame(source, destination string) (bool, error) {
	args := m.Called(source, destination)
	return args.Bool(0), args.Error(1)
}

func (m *MockFileUtils) ReadDirectory(path string) ([]os.DirEntry, error) {
	args := m.Called(path)
	return args.Get(0).([]os.DirEntry), args.Error(1)
}

func (m *MockFileUtils) DeleteAll(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

type MockDirEntry struct {
	mock.Mock
}

func (m *MockDirEntry) Name() string {
	args := m.Called()
	return args.String(0)
}

func (m *MockDirEntry) IsDir() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockDirEntry) Type() os.FileMode {
	args := m.Called()
	return args.Get(0).(os.FileMode)
}

func (m *MockDirEntry) Info() (os.FileInfo, error) {
	args := m.Called()
	return args.Get(0).(os.FileInfo), args.Error(1)
}
