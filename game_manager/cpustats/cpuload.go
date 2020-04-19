package cpustats

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sowa-gregory/game_cabinet/game_manager/utils/atomicint"
)

const (
	statPath          = "/proc/stat"
	statusLoadNew     = 0
	statusLoadStarted = 1
	statusLoadStopped = 2

	// DefaultFreq - default frequency of cpu statistics collection
	DefaultFreq = 2
)

// CPUStats - structure representing internal state of the object
type CPUStats struct {
	status  *atomicint.AtomicInt
	channel chan map[string]uint
	mutex   sync.Mutex
	wait    sync.WaitGroup
}

var instance *CPUStats
var once sync.Once

type cpuInfo struct {
	cpuID       string
	total, idle uint64
}

// GetInstance - creates singleton instance of CPULoad
func GetInstance() *CPUStats {
	once.Do(func() {
		var cpuStats CPUStats
		cpuStats.status = atomicint.New(statusLoadNew)
		instance = &cpuStats
	})
	return instance
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

// StopLoadMeasure - stops collection of cpu statistics. Gracefully finished backgoud goroutine.
func (cpuLoadObj *CPUStats) StopLoadMeasure() {
	cpuLoadObj.mutex.Lock()
	defer cpuLoadObj.mutex.Unlock()

	if cpuLoadObj.status.Load() == statusLoadStarted {

		cpuLoadObj.status.Store(statusLoadStopped)
		for range cpuLoadObj.channel {
		}
		cpuLoadObj.wait.Wait()
		print("stopload measure waited")
		cpuLoadObj.status.Store(statusLoadNew)
	} else {
		log.Panic("cpuload can be stopped only when started")
	}
}

func (cpuLoadObj *CPUStats) loadRoutine(frequency int) {
	defer cpuLoadObj.wait.Done()

	for {
		startInfo, err := getCPUInfo()
		if err != nil {
			log.Print("cpuload ", err)
			continue
		}

		if cpuLoadObj.status.Load() == statusLoadStopped {
			close(cpuLoadObj.channel)
			return
		}
		time.Sleep(time.Duration(frequency) * time.Second)
		if cpuLoadObj.status.Load() == statusLoadStopped {
			close(cpuLoadObj.channel)
			return
		}

		endInfo, err := getCPUInfo()
		if err != nil {
			log.Print("cpuload ", err)
			continue
		}

		load := make(map[string]uint)
		for index := 0; index < len(startInfo); index++ {
			total := (endInfo[index].total - startInfo[index].total)
			idle := (endInfo[index].idle - startInfo[index].idle)
			cpuPercentage := (uint)(100 * (total - idle) / total)
			load[startInfo[index].cpuID] = cpuPercentage
		}
		cpuLoadObj.channel <- load
	}
}

// StartLoadMeasure - starts collection of cpu statistics, which are send to channel in background
func (cpuLoadObj *CPUStats) StartLoadMeasure(frequency int) chan map[string]uint {
	cpuLoadObj.mutex.Lock()
	defer cpuLoadObj.mutex.Unlock()

	if cpuLoadObj.status.Load() != statusLoadNew {
		log.Panic("cpuload already started")
	}
	cpuLoadObj.status.Store(statusLoadStarted)
	cpuLoadObj.channel = make(chan map[string]uint)
	cpuLoadObj.wait.Add(1)
	go cpuLoadObj.loadRoutine(frequency)
	return cpuLoadObj.channel
}
