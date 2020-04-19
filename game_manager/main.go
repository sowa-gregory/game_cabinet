package main

import (
	"fmt"
	"time"

	"github.com/sowa-gregory/game_cabinet/game_manager/cpuinfo"
)

func main() {

	c := cpuinfo.GetLoad(2)

	d := time.After(2 * time.Second)

	a := cpuinfo.GetTemp()
	print(a)
	for {
		select {
		case load := <-c:
			fmt.Println(load)
			c = cpuinfo.GetLoad(2)

		case <-d:
			fmt.Println("timer")
			d = time.After(5 * time.Second)
		}
	}
}
