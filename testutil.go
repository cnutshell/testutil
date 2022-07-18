package testutil

import (
	"sync"
	"testing"
)

type DnStore interface {
	Start()
	Stop()

	ID() string
}

type LogStore interface {
	Start()
	Stop()

	ID() string
}

// FIXME: verification utilities
// FIXME: integrate with real log store
// FIXME: integrate with real dn store
type TestCluster interface {
	// Start starts stores sequentially
	Start()
	// Stop stops stores sequentially
	Stop()

	// RegisterDnStore registers dn stores
	RegisterDnStore(stores ...DnStore) error
	// GetDnStore fetches dn store instance
	GetDnStore(storeID string) (DnStore, error)

	// RegisterLogStore registers log stores
	RegisterLogStore(stores ...LogStore) error
	// GetLogStore fetches log store instance
	GetLogStore(storeID string) (LogStore, error)

	/*
		StartNetworkPartition(partitions [][]int)
		StopNetworkPartition()
	*/
}

func NewTestCluster(t *testing.T, dnStores []DnStore, logStores []LogStore) (TestCluster, error) {
	c := &testCluster{}

	c.dn.stores = make(map[string]DnStore)
	c.log.stores = make(map[string]LogStore)

	if err := c.registerDnStores(dnStores...); err != nil {
		return nil, err
	}

	if err := c.registerLogStores(logStores...); err != nil {
		return nil, err
	}

	return c, nil
}

type testCluster struct {
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

func (t *testCluster) registerDnStores(stores ...DnStore) error {
	t.dn.Lock()
	defer t.dn.Unlock()

	for _, s := range stores {
		id := s.ID()
		if _, ok := t.dn.stores[id]; ok {
			return wrappedError(ErrStoreDuplicated, id)
		}
		t.dn.stores[id] = s
	}

	return nil
}

func (t *testCluster) registerLogStores(stores ...LogStore) error {
	t.log.Lock()
	defer t.log.Unlock()

	for _, s := range stores {
		id := s.ID()
		if _, ok := t.log.stores[id]; ok {
			return wrappedError(ErrStoreDuplicated, id)
		}
		t.log.stores[id] = s
	}

	return nil
}

func (t *testCluster) Start() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.mu.running {
		return
	}

	for _, ds := range t.dn.stores {
		ds.Start()
	}

	for _, ls := range t.log.stores {
		ls.Start()
	}

	t.mu.running = true
}

func (t *testCluster) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.mu.running {
		return
	}

	for _, ds := range t.dn.stores {
		ds.Stop()
	}

	for _, ls := range t.log.stores {
		ls.Stop()
	}

	t.mu.running = false
}

func (t *testCluster) RegisterDnStore(stores ...DnStore) error {
	return t.registerDnStores(stores...)
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

func (t *testCluster) RegisterLogStore(stores ...LogStore) error {
	return t.registerLogStores(stores...)
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
