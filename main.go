package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"sync"
)

// Declaration of global constants
const (
	BLOCK_SIZE = 128 * 1024 // 128 KiB
	DICT_SIZE  = 32 * 1024  // 32 KiB
	SUM_SIZE   = 8          // 8 bytes
)

// Parsing processes flag
var processes int

func init() {
	var (
		defaultProcesses = runtime.NumCPU()
		usage            = "Specify number of goroutines to use for compression"
	)
	flag.IntVar(&processes, "processes", defaultProcesses, usage)
	flag.IntVar(&processes, "p", defaultProcesses, usage)
}

// This implementation of concurrent compression utilizes the pipelined,
// fan-out, fan-in concurrency pattern as described in
// https://go.dev/blog/pipelines
// There are three stages for a Block to be pipelined through:
// (1) Read stage
// (2) Compress stage
// (3) Write stage
// Just realized, we want compress to happen concurrent with read...
// This current strategy doesn't do that... We'll benchmark later.
func main() {
	flag.Parse()

	r := read()

	c := compress(r)

	for output := range c {
		write(output)
	}

	/*
		compressOutbounds := make([]<-chan *block, processes)
		for p := 0; p < processes; p++ {
			compressOutbounds[p] = compress(r)
		}

		for c := range mergeSlice(compressOutbounds) {
			write(c)
		}
	*/
}

// Read stage
func read() <-chan *block {
	log.Println("reading")
	out := make(chan *block)

	go func() {
		// Start reading input from Stdin in byte array buffers with BLOCK_SIZE
		inputBuffer := make([]byte, BLOCK_SIZE, BLOCK_SIZE)

		var numBytesTotal int
		var numBlocks int

		reader := bufio.NewReader(os.Stdin)
		numBytes, err := reader.Read(inputBuffer)
		for err != io.EOF {
			numBytesTotal += numBytes
			numBytes, err = reader.Read(inputBuffer)

			// check if readBuffer is the last block in the buffer
			isLastBlock := false
			_, err = reader.Peek(1)
			if err == bufio.ErrNegativeCount {
				isLastBlock = true
			}

			numBlocks++

			b := block{
				Index:     numBlocks,
				LastBlock: isLastBlock,
				Data:      inputBuffer,
			}
			log.Println("read block#" + strconv.Itoa(b.Index))
			out <- &b
		}
	}()

	return out
}

// Compress stage
func compress(in <-chan *block) <-chan *block {
	out := make(chan *block)

	go func() {
		for b := range in {
			log.Println("compressing block#" + strconv.Itoa(b.Index))
			out <- b
		}
		close(out)
	}()

	return out
}

// Write stage
func write(b *block) {
	log.Println("writing block#" + strconv.Itoa(b.Index))
}

// mergeList fans-in slice of results from the compress goroutines into the write stage
func mergeSlice(compressOutbounds []<-chan *block) <-chan *block {
	log.Println("merging")

	var wg sync.WaitGroup
	out := make(chan *block)

	wg.Add(len(compressOutbounds))
	for i := 0; i < len(compressOutbounds); i++ {
		go func() {
			for b := range compressOutbounds[i] {
				out <- b
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func merge(cs ...<-chan *block) <-chan *block {
	return nil
}

type block struct {
	Index     int
	LastBlock bool
	Data      []byte
	CRC32     []byte
	Err       error
}
