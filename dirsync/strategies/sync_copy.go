package strategies

import (
	"fmt"
	"path/filepath"
	"sync"

	fileutils "github.com/Matyjash/DirSync/dirsync/fileutils"
)

const (
	fileErrorMsg = "file: %s, skipping due to ERROR: %v\n"

	defaultMaxConcurrentFiles = 5
)

type syncCopy struct {
	semaphore chan struct{}
	wg        sync.WaitGroup
	fileUtils fileutils.FileUtils
}

func NewSyncCopy() *syncCopy {
	return &syncCopy{
		semaphore: make(chan struct{}, defaultMaxConcurrentFiles),
		wg:        sync.WaitGroup{},
		fileUtils: fileutils.NewDefaultFileHandler(),
	}
}

func (sc *syncCopy) SyncDirectories(source, destination string) {
	sc.syncOverDirectory(source, destination)
	sc.wg.Wait()
}

func (sc *syncCopy) syncOverDirectory(source, destination string) {
	files, err := sc.fileUtils.ReadDirectory(source)
	if err != nil {
		fmt.Printf("directory: %s, failed to read contents due to ERROR: %v\n", source, err)
		return
	}

	for _, file := range files {
		sourceFilePath := filepath.Join(source, file.Name())
		destinationFilePath := filepath.Join(destination, file.Name())
		if file.IsDir() {
			sc.syncDirectory(sourceFilePath, destinationFilePath)
		} else {
			sc.wg.Add(1)
			sc.semaphore <- struct{}{}
			go func(sourceFile, destinationFile string) {
				defer sc.wg.Done()
				defer func() {
					<-sc.semaphore
				}()
				sc.syncFile(sourceFile, destinationFile)
			}(sourceFilePath, destinationFilePath)
		}
	}
}

func (sc *syncCopy) syncFile(source, destination string) {
	exists, err := sc.fileUtils.Exists(destination)
	if err != nil {
		fmt.Printf(fileErrorMsg, source, err)
		return
	}
	if exists {
		same, err := sc.fileUtils.IsTheSame(source, destination)
		if err != nil {
			fmt.Printf(fileErrorMsg, source, err)
			return
		}
		if same {
			fmt.Printf("file: %s, already exists and is the same, skipping\n", source)
			return
		}
	}

	err = sc.fileUtils.CopyFile(source, destination)
	if err != nil {
		fmt.Printf("file: %s, failed to copy to destination due to ERROR: %v\n", source, err)
		return
	}
	fmt.Printf("file: %s, copied to destination\n", source)
}

func (sc *syncCopy) syncDirectory(source, destination string) {
	exists, err := sc.fileUtils.Exists(destination)
	if err != nil {
		fmt.Printf("directory: %s, skipping due to ERROR: %v\n", source, err)
		return
	}
	if !exists {
		err = sc.fileUtils.RecreateDirectory(source, destination)
		if err != nil {
			fmt.Printf("directory: %s, failed to create in destination due to ERROR: %v\n", source, err)
			return
		}
		fmt.Printf("directory: %s, created in destination\n", source)
	}
	sc.syncOverDirectory(source, destination)
}
