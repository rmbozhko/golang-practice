package main

import (
	"context"
	"fmt"
	"math/rand/v2"
)

// Generator that returns random word from predefined slice.
func wordsGenerator() func() string {
	words := []string{"go", "is", "great", "go", "is", "fast"}
	return func() string {
		return words[rand.IntN(len(words))]
	}
}

// Read n-words from input
func readNWords(ctx context.Context, wordsCount int) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		defer fmt.Println("Done - reading words")
		generateWords := wordsGenerator()
		for range wordsCount {
			select {
			case out <- generateWords():
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// Filter only unique words (each word should appear once)
func keepUnique(ctx context.Context, in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		defer fmt.Println("Done - keeping unique words")
		dict := make(map[string]struct{})
		for word := range in {
			if _, seen := dict[word]; !seen {
				dict[word] = struct{}{}
				select {
				case out <- word:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}

func reverseSwapRunes(s string) string {
	chars := []rune(s)
	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}
	return string(chars)
}

// Flip each word (reverse the characters)
func reverseWord(ctx context.Context, in <-chan string) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		defer fmt.Println("Done - reversing words")
		for word := range in {
			select {
			case out <- reverseSwapRunes(word):
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func Pipeline() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	out := reverseWord(ctx,
		keepUnique(ctx,
			readNWords(ctx, 20)))

	// Send results to the consumer.
	for word := range out {
		fmt.Println(word)
		if word == "og" {
			fmt.Println("Found 'og', stopping...")
			break
		}
	}
}

func main() {
	Pipeline()
}
