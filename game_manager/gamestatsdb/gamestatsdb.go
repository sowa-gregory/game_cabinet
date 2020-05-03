package gamestatsdb

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/sowa-gregory/game_cabinet/game_manager/utils/routinestop"
)

type gameStatDB struct {
	dataChannel <-chan string
	regEx       *regexp.Regexp
	rtnStop     *routinestop.RoutineStop
	file        *os.File
	fileName    string
}

type gameStat struct {
	Time     string
	Event    string
	GameID   string
	Console  string
	GameName string
}

func New(dataChannel <-chan string) *gameStatDB {
	var instance gameStatDB
	instance.dataChannel = dataChannel
	instance.rtnStop = routinestop.New()
	instance.regEx = regexp.MustCompile(`^(START|END) (\S*) \"(.*)\"$`)
	instance.file = nil
	instance.fileName = ""
	return &instance
}

// example START snes "Super RType"

func (gameStatDBObj *gameStatDB) statLineToJson(statLine string) (string, error) {
	match := gameStatDBObj.regEx.FindStringSubmatch(statLine)
	if match == nil {
		return "", fmt.Errorf("invalid game stat line:%s", statLine)
	}
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	shaOut := sha256.Sum256(([]byte)(match[2] + " " + match[3]))

	stat := gameStat{timeStr, match[1], hex.EncodeToString(shaOut[:]), match[2], match[3]}
	jsonStat, err := json.Marshal(stat)
	if err != nil {
		return "", err
	}
	return string(jsonStat), err
}

func (gameStatDBObj *gameStatDB) processData(statLine string) {
	file, err := os.OpenFile(fileName, os.O_APPEND, 0644)
}

func (gameStatDBObj *gameStatDB) StartProcessing() {
	go func() {
		defer gameStatDBObj.rtnStop.Done()
		for {
			var statLine string
			select {
			case statLine = <-gameStatDBObj.dataChannel:
				gameStatDBObj.processData(statLine)
			case <-gameStatDBObj.rtnStop.GetStopChannel():
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

/*
func Test() {
	time := time.Now()
	g := gameStat{time.Format("2006-01-02 15:04:05"), 89, "strzelanka", "aad33"}
	out, _ := json.Marshal(g)
	fmt.Println(string(out))
}*/
