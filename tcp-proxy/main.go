package main

import (
	"fmt"
	"os"
)

type FooReader struct{}

func (foo *FooReader) Read(b []byte) (n int, err error) {
	fmt.Print("in > ")
	return os.Stdin.Read(b)
}

type FooWriter struct{}

func (foo *FooWriter) Write(b []byte) (n int, err error) {
	fmt.Print("out > ")
	return os.Stdout.Write(b)
}

func main() {
	// intantiate reader and writer
	var (
		reader FooReader
		writer FooWriter
	)

	input := make([]byte, 4096)
	s, err := reader.Read(input)
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("Read %d bytes from StdIn\n", s)

	w, err := writer.Write(input)
	if err != nil {
		os.Exit(1)
	}
	fmt.Printf("Write %d bytes to StdOut\n", w)
}
