package routinestop

import (
	"sync"

	"github.com/sowa-gregory/game_cabinet/game_manager/utils/atomicbool"
)

type RoutineStop struct {
	stopped     atomicbool.AtomicBool
	wait        *sync.WaitGroup
	StopChannel chan bool
}

func New() *RoutineStop {
	instance := RoutineStop{}
	instance.StopChannel = make(chan bool, 1)
	return &instance
}

func (routineStopObj *RoutineStop) RequestStop(waitGrp *sync.WaitGroup) bool {
	if routineStopObj.stopped.SwapIfFalse() {
		routineStopObj.wait = waitGrp
		routineStopObj.wait.Add(1)
		routineStopObj.StopChannel <- true
		return true
	}
	return false
}

func (routineStopObj *RoutineStop) Done() {
	routineStopObj.wait.Done()
}
