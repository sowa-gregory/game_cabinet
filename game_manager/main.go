package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/sowa-gregory/game_cabinet/game_manager/cpuinfo"
)

func main() {

	c := cpuinfo.GetLoad()

	d := time.After(2 * time.Second)
	e := cpuinfo.GetTemperature()
	print(e)
	for {
		fmt.Println("gp", runtime.NumGoroutine())
		select {
	
		case load := <-c:
			fmt.Println(load)
			c = cpuinfo.GetLoad()

		case <-d:
			fmt.Println("timer")
			d = time.After(5 * time.Second)
		}
	}
}
