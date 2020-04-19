package atomicbool

import (
	"sync/atomic"
)

// AtomicBool class
type AtomicBool struct {
	value int32
}

const (
	stateFalse = 0
	stateTrue  = 1
)

// New - creates instance of AtomicInt
func New() *AtomicBool {
	instance := AtomicBool{}
	return &instance
}

// SwapIfFalse -
func (atomicBoolObj *AtomicBool) SwapIfFalse() bool {
	return atomic.CompareAndSwapInt32(&atomicBoolObj.value, stateFalse, stateTrue)
}

// SwapIfTrue -
func (atomicBoolObj *AtomicBool) SwapIfTrue() bool {
	return atomic.CompareAndSwapInt32(&atomicBoolObj.value, stateTrue, stateFalse)
}

// Load - atomic read of stored value
func (atomicBoolObj *AtomicBool) Load() bool {
	if atomic.LoadInt32(&atomicBoolObj.value) == stateTrue {
		return true
	}
	return false
}

// Store - atomic store of new value
func (atomicBoolObj *AtomicBool) Store(value bool) {
	newValue := stateTrue
	if !value {
		newValue = stateFalse
	}
	atomic.StoreInt32(&atomicBoolObj.value, int32(newValue))
}
