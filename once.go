package sync

import (
	"context"
	"sync/atomic"
	"time"
)

type Once interface {
	Do(f func() error) error
	DoTimeout(timeout time.Duration, f func() error) error
	DoContext(ctx context.Context, f func() error) error
}

// NewOnce returns a new instance of the Once interface.
func NewOnce() Once {
	return &once{
		m: NewMutex(),
	}
}

type once struct {
	m    Mutex
	done uint32
}

// Do executes the given function f if it has not been executed before.
// It returns nil if the function has already been executed.
func (o *once) Do(f func() error) error {
	if atomic.LoadUint32(&o.done) == 1 {
		return nil
	}

	o.m.Lock()
	defer o.m.Unlock()

	if o.done == 0 {
		return o.doContext(context.Background(), f)
	}

	return nil
}

// DoTimeout executes the given function f if it has not been executed before, with a specified timeout.
// It returns nil if the function has already been executed or if the timeout is <= 0.
// The function is executed in a goroutine and the execution is canceled if the timeout duration is reached.
func (o *once) DoTimeout(timeout time.Duration, f func() error) error {
	if atomic.LoadUint32(&o.done) == 1 {
		return nil
	}

	o.m.Lock()
	defer o.m.Unlock()

	if o.done == 0 {
		if timeout > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			return o.doContext(ctx, f)
		} else {
			return o.doContext(context.Background(), f)
		}
	}

	return nil
}

// DoContext executes the given function f if it has not been executed before, with the provided context ctx.
// It returns nil if the function has already been executed or if the context expires before the function completes.
func (o *once) DoContext(ctx context.Context, f func() error) error {
	if atomic.LoadUint32(&o.done) == 1 {
		return nil
	}

	o.m.Lock()
	defer o.m.Unlock()

	if o.done == 0 {
		return o.doContext(ctx, f)
	}

	return nil
}

func (o *once) doContext(ctx context.Context, f func() error) error {
	defer atomic.StoreUint32(&o.done, 1)

	doneCh := make(chan struct{})
	var funcErr error
	var panicErr any = nil
	go func() {
		defer close(doneCh)

		defer func() {
			if r := recover(); r != nil {
				panicErr = r
			}
		}()

		funcErr = f()
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-doneCh:
		if panicErr != nil {
			panic(panicErr)
		}

		return funcErr
	}
}
