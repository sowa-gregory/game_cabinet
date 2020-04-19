package cpuinfo

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	statPath = "/proc/stat"
	// DefaultFreq - default frequency of cpu statistics collection
	defaultTimeSpan = 2
)

type cpuInfo struct {
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

func getCPUInfo() ([]cpuInfo, error) {
	file, err := os.Open(statPath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	err = scanner.Err()
	if err != nil {
		return nil, err
	}

	var cpuInfoArr []cpuInfo
	for scanner.Scan() {
		line := strings.ReplaceAll(scanner.Text(), "  ", " ")
		if strings.HasPrefix(line, "cpu") {
			// finds end of cpu name - typically cpu, cpu0, cpu1....
			splitIndex := strings.IndexByte(line, ' ')
			cpuID := line[:splitIndex]
			statVals, err := statLineToUint(line[splitIndex+1:])
			if err != nil {
				return nil, err
			}
			// cpu_id user nice system idle iowait irq softrig steal
			idle := statVals[3] + statVals[4]                                                            // idle + iowait
			nonIdle := statVals[0] + statVals[1] + statVals[2] + statVals[5] + statVals[6] + statVals[7] // user+nice+system+irq+softrig+steal
			cpuInfoArr = append(cpuInfoArr, cpuInfo{cpuID, idle + nonIdle, idle})
		}
	}
	return cpuInfoArr, nil
}

func loadRoutine(channel chan LoadResult) {
	startInfo, err := getCPUInfo()
	if err != nil {
		channel <- LoadResult{nil, err}
		log.Print("cpuload ", err)
		return
	}

	time.Sleep(time.Duration(defaultTimeSpan) * time.Second)

	endInfo, err := getCPUInfo()
	if err != nil {
		channel <- LoadResult{nil, err}

		log.Print("cpuload ", err)
		return
	}

	load := make(map[string]uint)
	for index := 0; index < len(startInfo); index++ {
		total := (endInfo[index].total - startInfo[index].total)
		idle := (endInfo[index].idle - startInfo[index].idle)
		cpuPercentage := (uint)(100 * (total - idle) / total)
		load[startInfo[index].cpuID] = cpuPercentage
	}
	channel <- LoadResult{load, nil}
}

type LoadResult struct {
	Load map[string]uint
	Err  error
}

// StartLoadMeasure - starts collection of cpu statistics, which are send to channel in background
func GetLoad() chan LoadResult {
	channel := make(chan LoadResult)
	go loadRoutine(channel)
	return channel
}

func GetTemperature() chan uint {

	channel := make(chan uint)

	go func() {
		content, err := ioutil.ReadFile("/tmp/a")
		if err != nil {
			log.Panicln("GetTemperature", err)
		}
		str := strings.ReplaceAll(string(content), "\n", "")
		temp, err := strconv.ParseUint(str, 10, 32)
		if err != nil {
			log.Panicln("GetTemperature parsing value", err)
		}
		channel <- uint(temp / 1000)
	}()
	return channel
}
