package sync

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeout(t *testing.T) {
	once := NewOnce()

	err := once.DoTimeout(100*time.Millisecond, func() error {
		time.Sleep(200 * time.Millisecond)
		return nil
	})

	assert.Equal(t, context.DeadlineExceeded, err)
}

type one int

func (o *one) Increment() {
	*o++
}

func run(t *testing.T, once Once, o *one, c chan bool) {
	once.Do(func() error {
		o.Increment()
		return nil
	})
	if v := *o; v != 1 {
		t.Errorf("once failed inside run: %d is not 1", v)
	}
	c <- true
}

func TestOnce(t *testing.T) {
	o := new(one)
	once := NewOnce()
	c := make(chan bool)
	const N = 10
	for i := 0; i < N; i++ {
		go run(t, once, o, c)
	}
	for i := 0; i < N; i++ {
		<-c
	}
	if *o != 1 {
		t.Errorf("once failed outside run: %d is not 1", *o)
	}
}

func TestOncePanic(t *testing.T) {
	once := NewOnce()

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Fatalf("once.Do did not panic")
			}
		}()
		once.Do(func() error {
			panic("failed")
		})
	}()

	once.Do(func() error {
		t.Fatalf("once.Do called twice")
		return nil
	})
}

func BenchmarkOnce(b *testing.B) {
	once := NewOnce()

	f := func() error { return nil }
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			once.Do(f)
		}
	})
}
