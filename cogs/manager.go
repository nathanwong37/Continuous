package cogs

//Manager struct is used in order to deal with Workers/Timers within the Shard
type Manager struct {
	shardID int
	worker  []Worker
}

//NewManger is a constructor for Manager
func NewManager(ID int) *Manager {
	return &Manager{shardID: ID, worker: make([]Worker, 0)}
}
