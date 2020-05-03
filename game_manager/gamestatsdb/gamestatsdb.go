package gamestatsdb

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/sowa-gregory/game_cabinet/game_manager/utils/routinestop"
)

type gameStatDB struct {
	dataChannel <-chan string
	rtnStop     *routinestop.RoutineStop
}

type gameStat struct {
	Time      string
	EventType int
	GameName  string
	GameID    string
}

func New(dataChannel <-chan string) *gameStatDB {
	var instance gameStatDB
	instance.dataChannel = dataChannel
	instance.rtnStop = routinestop.New()
	return &instance
}

func (gameStatDBObj *gameStatDB) processData(data string) {

}

func (gameStatDBObj *gameStatDB) StartProcessing() {
	go func() {
		defer fmt.Println("exit db")
		for {
			var data string
			select {
			case data = <-gameStatDBObj.dataChannel:
				fmt.Println("stat:", data)
				gameStatDBObj.processData(data)
			case <-gameStatDBObj.rtnStop.GetStopChannel():
				gameStatDBObj.rtnStop.Done()
				return
			}
		}
	}()
}

func (gameStatDBObj *gameStatDB) Stop(waitGrp *sync.WaitGroup) {
	if !gameStatDBObj.rtnStop.RequestStop(waitGrp) {
		log.Panicln("GameStatDB already stopped")
	}

}

func Test() {
	time := time.Now()
	g := gameStat{time.Format("2006-01-02 15:04:05"), 89, "strzelanka", "aad33"}
	out, _ := json.Marshal(g)
	fmt.Println(string(out))
}
