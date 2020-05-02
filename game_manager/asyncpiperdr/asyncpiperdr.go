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
	channel chan string
	wait    *sync.WaitGroup
	stopped atomicbool.AtomicBool
}

// New - creates instance of asyncPipeRdr
func New() *asyncPipeRdr {
	var instance asyncPipeRdr
	instance.channel = make(chan string, 1) // must have length 1 to buffer write otherwise write to stopchannel may block
	instance.wait = &sync.WaitGroup{}
	return &instance
}

func processRead(pipe string, channel chan string) (continuation bool) {
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
			case <-channel:
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
		if by == '\n' {
			select {
			case channel <- string(buffer):
			case <-channel:
				return false
			}
			buffer = nil
		} else {
			buffer = append(buffer, by)
		}
	}
}

func (asyncPipeRdrObj *asyncPipeRdr) Stop() {
	if asyncPipeRdrObj.stopped.SwapIfFalse() {
		asyncPipeRdrObj.wait.Add(1)
		asyncPipeRdrObj.channel <- "stop"
		asyncPipeRdrObj.wait.Wait()
	} else {
		log.Panicln("AsyncPipeRdr already stopped")
	}
}

func (asyncPipeRdrObj *asyncPipeRdr) IsPipeExists(pipe string) bool {
	if fileInfo, err := os.Stat(pipe); err == nil && fileInfo.Mode()&os.ModeNamedPipe >= 0 {
		return true
	}
	return false
}

func (asyncPipeRdrObj *asyncPipeRdr) Read(pipe string) <-chan string {
	fileInfo, err := os.Stat(pipe)
	if err == nil {
		if fileInfo.Mode()&os.ModeNamedPipe == 0 {
			log.Panicf("AsyncPipeRdr path: %s is not a named pipe\n", pipe)
		}
	} else {
		log.Panicf("AsyncPipeRdr path: %s doesn't exist\n", pipe)
	}

	go func() {
		defer fmt.Println("exit routing")
		for processRead(pipe, asyncPipeRdrObj.channel) {
		}
		close(asyncPipeRdrObj.channel)
		asyncPipeRdrObj.wait.Done()
	}()
	return asyncPipeRdrObj.channel
}
