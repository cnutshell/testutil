package testutil

const (
	defaultDnStoreCount  = 1
	defaultLogStoreCount = 1

	defaultDnShardCount  = 1
	defaultLogShardCount = 1

	defaultLogReplicaNum = 3
)

type Options struct {
	dnStoreCount  int
	logStoreCount int

	dnShardCount  uint64
	logShardCount uint64
	logReplicaNum uint64
}

func DefaultOptions() Options {
	return Options{
		dnStoreCount:  defaultDnStoreCount,
		logStoreCount: defaultLogStoreCount,
		dnShardCount:  defaultDnShardCount,
		logShardCount: defaultLogShardCount,
		logReplicaNum: defaultLogReplicaNum,
	}
}

func validateOptions(opt *Options) {
	if opt.dnStoreCount == 0 {
		opt.dnStoreCount = defaultDnStoreCount
	}
	if opt.logStoreCount == 0 {
		opt.logStoreCount = defaultLogStoreCount
	}
	if opt.dnShardCount == 0 {
		opt.dnShardCount = defaultDnShardCount
	}
	if opt.logShardCount == 0 {
		opt.logShardCount = defaultLogShardCount
	}
	if opt.logReplicaNum == 0 {
		opt.logReplicaNum = defaultLogReplicaNum
	}
}

func (opt Options) WithDnStoreCount(count int) Options {
	opt.dnStoreCount = count
	return opt
}

func (opt Options) WithLogStoreCount(count int) Options {
	opt.logStoreCount = count
	return opt
}

func (opt Options) WithLogShardCount(count uint64) Options {
	opt.logShardCount = count
	return opt
}

func (opt Options) WithDnShardCount(count uint64) Options {
	opt.dnShardCount = count
	return opt
}

func (opt Options) WithLogReplicaNum(num uint64) Options {
	opt.logReplicaNum = num
	return opt
}

// ......
