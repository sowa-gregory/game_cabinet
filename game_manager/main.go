package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/sowa-gregory/game_cabinet/game_manager/asyncpiperdr"
//	"github.com/sowa-gregory/game_cabinet/game_manager/cpuinfo"
	"github.com/sowa-gregory/game_cabinet/game_manager/gamestatsdb"
)

func dbtest() {
	gamestatsdb.Test()
}

func Test() {
	rdr := asyncpiperdr.New()

	c := rdr.StartReading("/tmp/gamestats")
	db := gamestatsdb.New(c)
	db.StartProcessing()

	var waitg sync.WaitGroup

	time.Sleep(time.Second * 10)
	db.Stop(&waitg)
	rdr.Stop(&waitg)

	waitg.Wait()
	fmt.Println("")
}

func main() {
	//dbtest()
	Test()
	os.Exit(1)

}
