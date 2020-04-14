package cpuload

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const stat = "/proc/stat"

type CPUStats struct {
	cpuID                                                 string
	user, nice, system, idle, iowait, irq, softrig, steal uint64
}

type CPUInfo struct {
	cpuID       string
	total, idle uint64
}

func statLineToUint(line string) ([]uint64, error) {
	var err error
	strArr := strings.Split(line, " ")
	outArr := make([]uint64, len(strArr))
	for ind, str := range strArr {
		outArr[ind], err = strconv.ParseUint(str, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	return outArr, nil
}

//GetStats aa
func GetStats() ([]CPUStats, error) {
	file, err := os.Open(stat)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	var cpuArr []CPUStats

	scanner := bufio.NewScanner(file)
	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	for scanner.Scan() {
		line := strings.ReplaceAll(scanner.Text(), "  ", " ")
		if strings.HasPrefix(line, "cpu") {
			var cpu CPUStats
			// finds end of cpu name - typically cpu, cpu0, cpu1....
			splitIndex := strings.IndexByte(line, ' ')

			cpu.cpuID = line[:splitIndex]
			statVals, err := statLineToUint(line[splitIndex+1:])
			if err != nil {
				return nil, err
			}

			// cpu_id user nice system idle iowait irq softrig steal
			cpu.user = statVals[1]
			cpu.nice = statVals[2]
			cpu.system = statVals[3]
			cpu.idle = statVals[4]
			cpu.iowait = statVals[5]
			cpu.irq = statVals[6]
			cpu.softrig = statVals[7]
			cpu.steal = statVals[8]

			cpuArr = append(cpuArr, cpu)
		}
	}

	return cpuArr, nil
}

// StartsToLoad converts actual stats to current
func StatsToCPUInfo(cpuStatsArr []CPUStats) []CPUInfo {
	cpuInfoArr := make([]CPUInfo, len(cpuStatsArr))
	for index, cpuStats := range cpuStatsArr {
		idle := cpuStats.idle + cpuStats.iowait
		nonIdle := cpuStats.user + cpuStats.nice + cpuStats.system + cpuStats.irq + cpuStats.softrig + cpuStats.steal
		total := idle + nonIdle
		cpuInfoArr[index] = CPUInfo{cpuStats.cpuID, total, idle}
	}
	return cpuInfoArr
}

func GetCPULoad() {
	statsArr, _ := GetStats()
	startInfo := StatsToCPUInfo(statsArr)
	time.Sleep(5 * time.Second)
	statsArr, _ = GetStats()
	endInfo := StatsToCPUInfo(statsArr)
	fmt.Println(startInfo)
	fmt.Println(endInfo)

	for index := 0; index < len(startInfo); index++ {
		total := (float32)(endInfo[index].total - startInfo[index].total)
		idle := (float32)(endInfo[index].idle - startInfo[index].idle)
		var cpuPercentage float32 = (total - idle) / total * 100.0
		fmt.Println(startInfo[index].cpuID)
		fmt.Println(cpuPercentage)
	}

}
