package asyncpiperdr

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/sowa-gregory/game_cabinet/game_manager/utils/atomicbool"
)

type asyncPipeRdr struct {
	// channel to send lines to consumers
	channel chan string
	// internal channel to indicate end of processing
	stopChannel chan bool
	wait        *sync.WaitGroup
	stopped     atomicbool.AtomicBool
}

// New - creates instance of asyncPipeRdr
func New() *asyncPipeRdr {
	var instance asyncPipeRdr
	instance.channel = make(chan string)
	instance.stopChannel = make(chan bool, 1) // must have length 1 to buffer write otherwise write to stopchannel may block
	instance.wait = nil
	return &instance
}

func processRead(pipe string, channel chan<- string, stpChan <-chan bool) (continuation bool) {
	file, err := os.OpenFile(pipe, os.O_RDWR, 0)
	defer file.Close()
	if err != nil {
		log.Panicln(err)
	}
	reader := bufio.NewReader(file)
	var buffer []byte

	for {
		// we need to read byte by byte - which allowes to read all data even if write to the pipe is really slow
		// usage of bufio.ReadString loses data if timeout occures during reading of line
		file.SetReadDeadline(time.Now().Add(1 * time.Second))
		by, err := reader.ReadByte()
		if _, ok := err.(*os.PathError); ok {
			// just timeout
			select {
			case <-stpChan:
				return false
			default:
			}
			continue
		}
		if err == io.EOF {
			return true
		}
		if err != nil {
			log.Println(err)
			return true
		}

		// allows for exit at reception of any byte from pipe
		select {
		case <-stpChan:
			return false
		default:
		}

		if by == '\n' {
			select {
			case channel <- string(buffer):
			case <-stpChan:
				return false
			}
			buffer = nil
		} else {
			buffer = append(buffer, by)
		}
	}
}

func (asyncPipeRdrObj *asyncPipeRdr) Stop(waitGrp *sync.WaitGroup) {
	if asyncPipeRdrObj.stopped.SwapIfFalse() {
		asyncPipeRdrObj.wait = waitGrp
		asyncPipeRdrObj.wait.Add(1)
		asyncPipeRdrObj.stopChannel <- true
	} else {
		log.Panicln("AsyncPipeRdr already stopped")
	}
}

func (asyncPipeRdrObj *asyncPipeRdr) PipeExists(pipe string) bool {
	if fileInfo, err := os.Stat(pipe); err == nil && fileInfo.Mode()&os.ModeNamedPipe >= 0 {
		return true
	}
	return false
}

func (asyncPipeRdrObj *asyncPipeRdr) StartReading(pipeName string) <-chan string {
	fileInfo, err := os.Stat(pipeName)
	if err == nil {
		if fileInfo.Mode()&os.ModeNamedPipe == 0 {
			log.Panicf("AsyncPipeRdr path: %s is not a named pipe\n", pipeName)
		}
	} else {
		log.Panicf("AsyncPipeRdr path: %s doesn't exist\n", pipeName)
	}

	go func() {
		defer fmt.Println("exit routing")
		for processRead(pipeName, asyncPipeRdrObj.channel, asyncPipeRdrObj.stopChannel) {
		}
		close(asyncPipeRdrObj.channel)
		asyncPipeRdrObj.wait.Done()
	}()
	return asyncPipeRdrObj.channel
}
