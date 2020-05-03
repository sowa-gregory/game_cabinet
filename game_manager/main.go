package main

import (
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/sowa-gregory/game_cabinet/game_manager/asyncpiperdr"
	"github.com/sowa-gregory/game_cabinet/game_manager/cpuinfo"
	"github.com/sowa-gregory/game_cabinet/game_manager/gamestatsdb"
)

func dbtest() {
	gamestatsdb.Test()
}

func Test() {
	rdr := asyncpiperdr.New()

	c := rdr.StartReading("/tmp/test")
	db := gamestatsdb.New(c)
	db.StartProcessing()

	var waitg sync.WaitGroup

	time.Sleep(time.Second * 1)
	db.Stop(&waitg)
	rdr.Stop(&waitg)

	waitg.Wait()
}

func main() {
	//dbtest()
	Test()
	os.Exit(1)
	c := cpuinfo.GetLoad()

	d := time.After(2 * time.Second)
	e := cpuinfo.GetTemperature()

	for {
		fmt.Println("gp", runtime.NumGoroutine())
		select {
		case temp := <-e:
			fmt.Println(temp)

		case load := <-c:
			fmt.Println(load)
			c = cpuinfo.GetLoad()

		case <-d:
			fmt.Println("timer")
			d = time.After(5 * time.Second)
		}
	}
}
