package testutil

import "fmt"

var (
	ErrStoreDuplicated  = fmt.Errorf("store id duplicated")
	ErrStoreNotExist    = fmt.Errorf("store not exist")
	ErrNoLeaderHakeeper = fmt.Errorf("no leader hakeeper")
)

func wrappedError(err error, msg string) error {
	return fmt.Errorf("%w: %s", err, msg)
}
