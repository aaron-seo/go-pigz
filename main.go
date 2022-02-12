package main

import (
	"bufio"
	"flag"
	"io"
	"os"
	"runtime"
	"sync"
)

const (
	BLOCK_SIZE = 128 * 1024 // 128 KiB
	DICT_SIZE  = 32 * 1024  // 32 KiB
)

// For parsing processes flag
var processes int

func init() {
	var (
		defaultProcesses = runtime.NumCPU()
		usage            = "Specify number of goroutines to use for compression"
	)
	flag.IntVar(&processes, "processes", defaultProcesses, usage)
	flag.IntVar(&processes, "p", defaultProcesses, usage)
}

// This implementation of concurrent compression utilizes the pipelined
// fan-out, fan-in concurrency pattern as described in
// https://go.dev/blog/pipelines
// There are three stages for a Block to be pipelined through:
// (1) Read stage
// (2) Compress stage
// (3) Write stage
func main() {
	// Parse arguments
	flag.Parse()

	r := read()

	// TODO should be p num of outbounds for compress
	c1 := compress(r)
	c2 := compress(r)

	for c := range merge(c1, c2) {
		write(c)
	}
}

func read() <-chan Block {
	out := make(chan Block)

	// Start reading input from Stdin in byte array buffers with BLOCK_SIZE
	readBuffer := make([]byte, BLOCK_SIZE, BLOCK_SIZE)

	var numBytesTotal int
	var numBlocks int

	reader := bufio.NewReader(os.Stdin)
	numBytes, err := reader.Read(readBuffer)
	for err != io.EOF {
		numBytesTotal += numBytes
		numBytes, err = reader.Read(readBuffer)

		// check if readBuffer is the lastBlock in buffer
		lastBlock := false
		_, err = reader.Peek(1)
		if err == bufio.ErrNegativeCount {
			lastBlock = true
		}

		numBlocks++
		out <- new(Block{
			index:     numBlocks,
			lastBlock: false,
			inputBuffer, readBuffer,
		})
	}
	return out
}
func compress(in <-chan Block) <-chan Block {
}

func write(in <-chan Block) {
}

func merge(cs ...<-chan Block) <-chan Block {
	var wg sync.WaitGroup

}

type Block struct {
	index        int
	lastBlock    bool
	inputBuffer  []byte
	outputBuffer []byte
}
