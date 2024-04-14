package sync

import (
	"github.com/sasha-s/go-deadlock"
	"os"
	"strconv"
	"time"
)

var (
	threshold   = 100 * time.Millisecond
	l           = &logger{}
	debug       = false
	useDeadlock = false
)

func init() {
	if b, err := strconv.ParseBool(os.Getenv("SYNC_DEBUG")); err == nil {
		l.enableDebug = b
		debug = b
		l.Debugf("Set lock debug to: %v", debug)
	}

	if n, _ := strconv.Atoi(os.Getenv("SYNC_LOCK_THRESHOLD")); n > 0 {
		threshold = time.Duration(n) * time.Millisecond
		l.Debugf("Set lock threshold at %v", threshold)
	}

	if b, err := strconv.ParseBool(os.Getenv("SYNC_USE_DEADLOCK")); err == nil {
		useDeadlock = b
		l.Debugf("Set useDeadlock to %v", useDeadlock)
	}

	if n, _ := strconv.Atoi(os.Getenv("SYNC_DEADLOCK_TIMEOUT")); n > 0 {
		deadlock.Opts.DeadlockTimeout = time.Duration(n) * time.Second
		l.Debugf("Enabling lock deadlocking at %v", deadlock.Opts.DeadlockTimeout)
		useDeadlock = true
	}
}
