package dirsync

import (
	"github.com/Matyjash/DirSync/dirsync/strategies"
)

type DirSync struct {
	Strategy strategies.SyncStrategy
}

func NewDirSync(deleteMissing bool) *DirSync {
	var strategy strategies.SyncStrategy
	if deleteMissing {
		strategy = strategies.NewSyncCopyDelete()
	} else {
		strategy = strategies.NewSyncCopy()
	}
	return &DirSync{Strategy: strategy}
}

func (ds *DirSync) Sync(source, destination string) {
	ds.Strategy.SyncDirectories(source, destination)
}
