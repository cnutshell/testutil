package testutil

import (
	"sync"
	"testing"
	"time"

	logpb "github.com/matrixorigin/matrixone/pkg/pb/logservice"

	"github.com/google/uuid"
)

type TestCluster interface {
	// Start starts stores sequentially
	Start() error
	// Stop stops stores sequentially
	Stop() error

	// TestOperation
	// TestAwareness
	// TestAssertState
	// TestWaitState
}

// TestOperation supports cluster operation
type TestOperation interface {
	StopDnStore(storeID string) error
	StartDnStore(storeID string) error

	StopLogStore(storeID string) error
	StartLogStore(storeID string) error

	/*
		StartNetworkPartition(partitions [][]int) error
		StopNetworkPartition() error
	*/
}

// TestAwareness provides cluster information.
type TestAwareness interface {
	// ListDnStores lists all dn stores
	ListDnStores() []string
	// ListLogStores lists all log stores
	ListLogStores() []string

	// GetDnStore fetches dn store instance
	GetDnStore(storeID string) (DnStore, error)
	// GetLogStore fetches log store instance
	GetLogStore(storeID string) (LogStore, error)
	// GetClusterState fetches current cluster state
	GetClusterState() (*logpb.CheckerState, error)
}

// TestAssertState asserts current cluster state.
type TestAssertState interface {
	AssertDnShardCount(t *testing.T, expected int)
	AssertDnReplicaCount(t *testing.T, shardID uint64, expected int)

	AssertLogShardCount(t *testing.T, expected int)
	AssertLogReplicaCount(t *testing.T, shardID uint64, expected int)

	AssertLeaderHakeeperState(t *testing.T, expeted logpb.HAKeeperState)
}

// TestWaitState waits cluster state until timeout.
type TestWaitState interface {
	WaitDnShardByCount(count int, timeout time.Duration)
	WaitDnReplicaByCount(shardID uint64, count int, timeout time.Duration)

	WaitLogShardByCount(count int, timeout time.Duration)
	WaitLogReplicaByCount(shardID uint64, count int, timeout time.Duration)
}

func NewTestCluster(t *testing.T, opt Options) (TestCluster, error) {
	validateOptions(&opt)

	c := &testCluster{
		opt: opt,
	}
	c.dn.stores = make(map[string]DnStore)
	c.log.stores = make(map[string]LogStore)
	c.bootstrap.errChan = make(chan error, 1)

	// construct dn stores
	for i := 0; i < opt.dnStoreCount; i++ {
		// FIXME: integrate with real dn store
		ds := newDnStore()
		id := ds.ID()
		if _, ok := c.dn.stores[id]; ok {
			return nil, wrappedError(ErrStoreDuplicated, id)
		}
		c.dn.stores[id] = ds
	}

	// construct log stores
	for i := 0; i < opt.logStoreCount; i++ {
		// FIXME: integrate with real log store
		ls := newLogStore()
		id := ls.ID()
		if _, ok := c.log.stores[id]; ok {
			return nil, wrappedError(ErrStoreDuplicated, id)
		}
		c.log.stores[id] = ls
	}
	return c, nil
}

type testCluster struct {
	opt Options

	bootstrap struct {
		sync.Once
		errChan chan error
	}

	dn struct {
		sync.Mutex
		stores map[string]DnStore
	}

	log struct {
		sync.Mutex
		stores map[string]LogStore
	}

	mu struct {
		sync.Mutex
		running bool
	}
}

func (t *testCluster) Start() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.mu.running {
		return nil
	}

	// start heartbeat for dn store
	for _, ds := range t.dn.stores {
		if err := ds.Start(); err != nil {
			return err
		}
	}

	// start hearbeat for log store
	for _, ls := range t.log.stores {
		if err := ls.Start(); err != nil {
			return err
		}
	}

	bootstrap := func() {
		var err error
		defer func() {
			t.bootstrap.errChan <- err
		}()

		var leader LogStore
		var term uint64
		for _, ls := range t.log.stores {
			if isLeader, currTerm, _ := ls.IsLeaderHakeeper(); isLeader {
				leader = ls
				term = currTerm
				break
			}
		}

		if leader == nil {
			err = ErrNoLeaderHakeeper
			return
		}

		// set cluster initialized state
		err = leader.SetInitialClusterInfo(
			t.opt.logShardCount,
			t.opt.dnShardCount,
			t.opt.logReplicaNum,
		)
		if err != nil {
			return
		}

		state, err := leader.GetClusterState()
		if err != nil {
			return
		}

		// FIXME: get rid of bootstrap
		leader.Bootstrap(term, state)
	}

	// do bootstrap only once
	t.bootstrap.Do(bootstrap)

	// check bootstrap error
	err := <-t.bootstrap.errChan
	if err != nil {
		return err
	}

	t.mu.running = true
	return nil
}

func (t *testCluster) Stop() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.mu.running {
		return nil
	}

	for _, ds := range t.dn.stores {
		if err := ds.Stop(); err != nil {
			return err
		}
	}

	for _, ls := range t.log.stores {
		if err := ls.Stop(); err != nil {
			return err
		}
	}

	t.mu.running = false
	return nil
}

func (t *testCluster) GetDnStore(storeID string) (DnStore, error) {
	t.dn.Lock()
	defer t.dn.Unlock()

	s, ok := t.dn.stores[storeID]
	if !ok {
		return nil, wrappedError(ErrStoreNotExist, storeID)
	}
	return s, nil
}

func (t *testCluster) GetLogStore(storeID string) (LogStore, error) {
	t.log.Lock()
	defer t.log.Unlock()

	s, ok := t.log.stores[storeID]
	if !ok {
		return nil, wrappedError(ErrStoreNotExist, storeID)
	}
	return s, nil
}

func (t *testCluster) ListDnStores() []string {
	ids := make([]string, 0, len(t.dn.stores))
	for id := range t.dn.stores {
		ids = append(ids, id)
	}
	return ids
}

func (t *testCluster) ListLogStores() []string {
	ids := make([]string, 0, len(t.log.stores))
	for id := range t.log.stores {
		ids = append(ids, id)
	}
	return ids
}

func newDnStore() DnStore {
	return &mockStore{}
}

func newLogStore() LogStore {
	return &mockStore{}
}

type mockStore struct{}

func (ms *mockStore) Start() error {
	return nil
}

func (ms *mockStore) Stop() error {
	return nil
}

func (ms *mockStore) ID() string {
	return uuid.New().String()
}

func (ms *mockStore) IsLeaderHakeeper() (bool, uint64, error) {
	return false, 0, nil
}

func (ms *mockStore) GetClusterState() (*logpb.CheckerState, error) {
	return &logpb.CheckerState{}, nil

}
func (ms *mockStore) SetInitialClusterInfo(
	numOfLogShards, numOfDNShards, numOfLogReplicas uint64,
) error {
	return nil
}

func (ms *mockStore) Bootstrap(term uint64, state *logpb.CheckerState) {
}
