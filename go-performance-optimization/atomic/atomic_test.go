package atomic

import (
	"sync"
	"sync/atomic"
	"testing"
)

type atomicCounter struct {
	i int32
}

func AtomicAddOne(c *atomicCounter) {
	atomic.AddInt32(&c.i, 1)
}

type mutexCounter struct {
	i int32
	sync.Mutex
}

func MutexAddOne(c *mutexCounter) {
	c.Lock()
	c.i++
	c.Unlock()
}

func BenchmarkAtomicAddOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := new(atomicCounter)
		AtomicAddOne(c)
	}
}

func BenchmarkMutexAddOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		c := new(mutexCounter)
		MutexAddOne(c)
	}
}
