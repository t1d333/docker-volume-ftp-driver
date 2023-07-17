package statemngr

type StateManager interface {
	SyncState() error
	SaveState() error
}
