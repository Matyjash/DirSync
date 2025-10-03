package strategies

import (
	"fmt"
	"path/filepath"
	"sync"

	fileutils "github.com/Matyjash/DirSync/dirsync/fileutils"
)

type SyncCopyDelete struct {
	wg        sync.WaitGroup
	fileUtils fileutils.FileUtils
}

func NewSyncCopyDelete() *SyncCopyDelete {
	return &SyncCopyDelete{
		wg:        sync.WaitGroup{},
		fileUtils: fileutils.NewDefaultFileHandler(),
	}
}

func (scd *SyncCopyDelete) SyncDirectories(source, destination string) {
	syncCopy := NewSyncCopy()

	scd.wg.Add(1)
	go func() {
		defer scd.wg.Done()
		syncCopy.SyncDirectories(source, destination)
	}()

	scd.wg.Add(1)
	go func() {
		defer scd.wg.Done()
		scd.deleteNotInSource(source, destination)
	}()

	scd.wg.Wait()
}

func (scd *SyncCopyDelete) deleteNotInSource(source, destination string) {
	files, err := scd.fileUtils.ReadDirectory(destination)
	if err != nil {
		fmt.Printf("directory: %s, failed to read contents due to ERROR: %v\n", destination, err)
		return
	}

	for _, file := range files {
		sourceFilePath := filepath.Join(source, file.Name())
		destinationFilePath := filepath.Join(destination, file.Name())
		exists, err := scd.fileUtils.Exists(sourceFilePath)
		if err != nil {
			fmt.Printf("file/directory: %s, skipping deletion check due to ERROR: %v\n", sourceFilePath, err)
			continue
		}

		if !exists {
			err = scd.fileUtils.DeleteAll(destinationFilePath)
			if err != nil {
				fmt.Printf("file/directory: %s, failed to delete from destination due to ERROR: %v\n", destinationFilePath, err)
				continue
			}
			fmt.Printf("file/directory: %s, deleted from destination\n", destinationFilePath)
		} else if file.IsDir() {
			scd.deleteNotInSource(sourceFilePath, destinationFilePath)
		}
	}
}
