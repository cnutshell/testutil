package testutil

import (
	"go.uber.org/zap/zapcore"
)

type Options struct {
	tmpDir        string
	dnStoreCount  int
	logStoreCount int
	logLevel      zapcore.Level
	// ......
}

func DefaultOptions(dir string) Options {
	return Options{
		tmpDir: dir,
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

func (opt Options) WithLogLevel(lvl zapcore.Level) Options {
	opt.logLevel = lvl
	return opt
}

// ......
