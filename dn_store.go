package testutil

type DnStore interface {
	// Start sends heartbeat and start to handle command.
	Start() error

	// Stop stops store
	Stop() error

	// ID returns uuid of store
	ID() string
}
