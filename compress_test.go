package main

import (
	"bytes"
	"compress/flate"
	"io"
	"log"
	"math/rand"
	"reflect"
	"testing"
)

// Test compressing a single buffer of data
func TestCompressSingle(t *testing.T) {
	// for got
	in := make(chan *block)

	randomData := make([]byte, BLOCK_SIZE)
	rand.Read(randomData)

	testBlock := block{
		Index:     0,
		LastBlock: true,
		RawData:   randomData,
	}

	go func() {
		in <- &testBlock
		close(in)
	}()

	outChannel := compress(in)
	got := (<-outChannel).CompressedData

	// for want
	var buffer bytes.Buffer

	w, err := flate.NewWriter(&buffer, flate.DefaultCompression)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := io.Copy(w, bytes.NewReader(randomData)); err != nil {
		log.Fatal(err)
	}
	if err := w.Close(); err != nil {
		log.Fatal(err)
	}

	want := buffer.Bytes()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got and want not equal")
	}
}

// TODO Test compressing multiple blocks of data
func TestCompressMultiple(t *testing.T) {
}
