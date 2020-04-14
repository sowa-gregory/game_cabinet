package main

import (
	"fmt"

	"github.com/sowa-gregory/game_cabinet/game_manager/cpuload"
)

func main() {
	cpuload.GetCPULoad()

	fmt.Println("start")
}
