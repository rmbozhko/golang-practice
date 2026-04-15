package main

import (
	"bufio"
	"bytes"
	"log"
)

// Exercise 7.1: Using the ideas from ByteCounter, implement counters for words and for lines.
// You will find bufio.ScanWords useful.

type ByteCounter int

// counting bytes
func (c *ByteCounter) Write(p []byte) (int, error) {
	*c += ByteCounter(len(p))
	return len(p), nil
}

func scannerHelper(p []byte, split bufio.SplitFunc) (int, error) {
	scanner := bufio.NewScanner(bytes.NewReader(p))
	scanner.Split(split)

	count := 0
	for scanner.Scan() {
		count++
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return count, nil
}

type WordCounter int

// counting words
func (c *WordCounter) Write(p []byte) (int, error) {
	count, err := scannerHelper(p, bufio.ScanWords)
	*c += WordCounter(count)
	return count, err
}

type LineCounter int

// counting lines
func (c *LineCounter) Write(p []byte) (int, error) {
	count, err := scannerHelper(p, bufio.ScanLines)
	*c += LineCounter(count)
	return count, err
}

func main() {
	var byteCount ByteCounter
	var wordCount WordCounter
	var lineCount LineCounter

	input := "Hello, world!\nThis is a test.\nCounting bytes, words, and lines."
	byteCount.Write([]byte(input))
	wordCount.Write([]byte(input))
	lineCount.Write([]byte(input))

	log.Printf("bytes=%d words=%d lines=%d", byteCount, wordCount, lineCount)
}
