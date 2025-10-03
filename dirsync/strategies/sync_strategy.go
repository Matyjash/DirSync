package strategies

type SyncStrategy interface {
	SyncDirectories(source, destination string)
}
