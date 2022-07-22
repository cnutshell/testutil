package testutil

import logpb "github.com/matrixorigin/matrixone/pkg/pb/logservice"

type LogStore interface {
	// Start sends heartbeat and start to handle command.
	Start() error

	// Stop stops store
	Stop() error

	// ID returns uuid of store
	ID() string

	// IsLeaderHakeeper checks hakeeper information.
	IsLeaderHakeeper() (bool, uint64, error)
	// GetClusterState returns cluster information from hakeeper leader.
	GetClusterState() (*logpb.CheckerState, error)
	// SetInitialClusterInfo sets cluster initialize state.
	SetInitialClusterInfo(numOfLogShards, numOfDNShards, numOfLogReplicas uint64) error
	// Bootstrap would bootstrap cluster according to its initialize state.
	Bootstrap(term uint64, state *logpb.CheckerState)
}
