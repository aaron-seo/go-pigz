package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"runtime"
)

const (
	BLOCK_SIZE = 128 * 1024 // 128 KiB
	DICT_SIZE  = 32 * 1024  // 32 KiB
)

var processes int

func init() {
	var (
		defaultProcesses = runtime.NumCPU()
		usage            = "Specify number of goroutines to use for compression"
	)
	flag.IntVar(&processes, "processes", defaultProcesses, usage)
	flag.IntVar(&processes, "p", defaultProcesses, usage)
}

func main() {
	// Parse arguments
	flag.Parse()

	// List of write jobs
	blockIndex := 0

	writeTasks

	// Launch write goroutine
	go write()

	// Start reading input from Stdin in byte array buffers with BLOCK_SIZE
	inputBuffer := make([]byte, BLOCK_SIZE, BLOCK_SIZE)
	var nBytesRead int
	reader := bufio.NewReader(os.Stdin)
	numBytes, err := reader.Read(inputBuffer)
	for err != io.EOF {
		numBytesTotal += numBytes
		numBytes, err = reader.Read(inputBuffer)
	}

}

type Block struct {
	Index       int
	LastBlock   bool
	InputBuffer []byte
}

func compress(block Block) {
}

func write() {
	// Write header

	// Listen for write tasks in the pool

	// Write trailer
}
