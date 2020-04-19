package main

import (
	"fmt"
	"time"

	"github.com/sowa-gregory/game_cabinet/game_manager/cpustats"
)

func main() {

	cpu := cpustats.GetInstance()

	c := cpu.StartLoadMeasure(cpustats.DefaultFreq)
	cpu.StopLoadMeasure()

	time.Sleep(2 * time.Second)
	c = cpu.StartLoadMeasure(cpustats.DefaultFreq)
	for i := 0; i < 3; i++ {
		load, err := <-c
		fmt.Println(i, load, err)
	}
	cpu.StopLoadMeasure()
}
