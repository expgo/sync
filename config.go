package sync

import (
	"github.com/sasha-s/go-deadlock"
	"os"
	"strconv"
	"time"
)

var (
	logSlowLock       = false
	slowLockThreshold = 100 * time.Millisecond
	l                 = &logger{}
	useDeadlock       = false
)

func init() {
	if b, err := strconv.ParseBool(os.Getenv("SYNC_LOG_SLOW_LOCK")); err == nil {
		l.enableDebug = b
		logSlowLock = b
		l.Debugf("Set lock logSlowLock to: %v", logSlowLock)
	}

	if n, _ := strconv.Atoi(os.Getenv("SYNC_SLOW_LOCK_THRESHOLD")); n > 0 {
		slowLockThreshold = time.Duration(n) * time.Millisecond
		l.Debugf("Set lock slowLockThreshold at %v", slowLockThreshold)
		l.enableDebug = true
		logSlowLock = true
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
