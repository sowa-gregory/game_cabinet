package asyncpiperdr

import (
	"bufio"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/sowa-gregory/game_cabinet/game_manager/utils/routinestop"
)

type asyncPipeRdr struct {
	// channel to send lines to consumers
	channel chan string
	rtnStop *routinestop.RoutineStop
}

// New - creates instance of asyncPipeRdr
func New() *asyncPipeRdr {
	var instance asyncPipeRdr
	instance.channel = make(chan string)
	instance.rtnStop = routinestop.New()
	return &instance
}

func (asyncPipeRdrObj *asyncPipeRdr) processRead(pipe string) (continuation bool) {
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
			case <-asyncPipeRdrObj.rtnStop.GetStopChannel():
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
		case <-asyncPipeRdrObj.rtnStop.GetStopChannel():
			return false
		default:
		}

		if by == '\n' {
			select {
			case asyncPipeRdrObj.channel <- string(buffer):
			case <-asyncPipeRdrObj.rtnStop.GetStopChannel():
				return false
			}
			buffer = nil
		} else {
			buffer = append(buffer, by)
		}
	}
}

func (asyncPipeRdrObj *asyncPipeRdr) Stop(waitGrp *sync.WaitGroup) {
	if !asyncPipeRdrObj.rtnStop.RequestStop(waitGrp) {
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
		defer asyncPipeRdrObj.rtnStop.Done()
		for asyncPipeRdrObj.processRead(pipeName) {
		}
		close(asyncPipeRdrObj.channel)
	}()
	return asyncPipeRdrObj.channel
}
