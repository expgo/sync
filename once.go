package sync

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type Once struct {
	m    sync.Mutex
	done uint32
}

func (o *Once) Do(f func() error) error {
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

func (o *Once) DoTimeout(timeout time.Duration, f func() error) error {
	if atomic.LoadUint32(&o.done) == 1 {
		return nil
	}

	o.m.Lock()
	defer o.m.Unlock()

	if o.done == 0 {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		return o.doContext(ctx, f)
	}

	return nil
}

func (o *Once) DoContext(ctx context.Context, f func() error) error {
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

func (o *Once) doContext(ctx context.Context, f func() error) error {
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
