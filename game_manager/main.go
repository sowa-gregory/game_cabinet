package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/sowa-gregory/game_cabinet/game_manager/asyncpiperdr"
	"github.com/sowa-gregory/game_cabinet/game_manager/cpuinfo"
)

func Test() {
	rdr := asyncpiperdr.New()

	c := rdr.Read("/tmp/test")

	for j := range c {
		fmt.Println(j)
	}

}

func main() {
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
