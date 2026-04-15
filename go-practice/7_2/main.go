package main

import (
	"bytes"
	"io"
	"log"
)

// Exercise 7.2: Write a function CountingWriter with the signature below that, given an
// io.Writer, returns a new Writer that wraps the original, and a pointer to an int64 variable
// that at any moment contains the number of bytes written to the new Writer

type ByteCounter struct {
	count int64
}

func (c *ByteCounter) Write(p []byte) (int, error) {
	c.count += int64(len(p))
	return len(p), nil
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	var counter ByteCounter
	return io.MultiWriter(w, &counter), &counter.count
}

func main() {
	// Example usage
	var buf bytes.Buffer
	writer, count := CountingWriter(&buf)
	writer.Write([]byte("Hello, world!"))
	log.Printf("Bytes written: %d", *count)

	writer.Write([]byte(" Another line."))
	log.Printf("Bytes written: %d", *count)
}
