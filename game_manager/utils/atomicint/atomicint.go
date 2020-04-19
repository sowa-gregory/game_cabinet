package atomicint

import (
	"sync/atomic"
)

// AtomicInt class
type AtomicInt struct {
	value int32
}

// New - creates instance of AtomicInt
func New(initialValue int32) *AtomicInt {
	instance := AtomicInt{}
	return &instance
}

// CompareAndSwap - atomic operation compare and swap. On success returns true.
func (atomicIntObj *AtomicInt) CompareAndSwap(oldValue int32, newValue int32) bool {
	return atomic.CompareAndSwapInt32(&atomicIntObj.value, oldValue, newValue)
}

// Load - atomic read of stored value
func (atomicIntObj *AtomicInt) Load() int32 {
	return atomic.LoadInt32(&atomicIntObj.value)
}

// Store - atomic store of new value
func (atomicIntObj *AtomicInt) Store(newValue int32) {
	atomic.StoreInt32(&atomicIntObj.value, newValue)
}
