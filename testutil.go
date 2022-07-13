package testutil

import "testing"

type DnStore interface {
	Start()
	Stop()
	Restart()
	// ...
}

type LogStore interface {
	Start()
	Stop()
	Restart()
	// ...
}

type TestCluster interface {
	Start()
	Stop()
	Restart()

	StartNetworkPartition(partitions [][]int)
	StopNetworkPartition()

	GetDnStore(id string) DnStore
	GetLogStore(id string) LogStore
}

func NewTestCluster(t *testing.T) TestCluster {
	return nil
}

type testCluster struct {
}

// FIXME: implement `TestCluster`
func (t *testCluster) Start() {
}

// ...
