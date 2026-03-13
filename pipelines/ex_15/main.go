package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type Job struct {
	Path string
}

type Result struct {
	Path  string
	Words int
	Err   error
}

func generator(ctx context.Context, paths []string) <-chan Job {
	out := make(chan Job)
	go func() {
		defer close(out)
		for _, path := range paths {
			select {
			case out <- Job{Path: path}:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

func processFile(path string) (int, error) {
	paths := map[string]int{
		"file_1.txt": 15,
		"file_2.txt": 30,
		"file_7.txt": 45,
		"file_3.txt": 60,
		"file_4.txt": 60,
	}
	if path == "file_7.txt" {
		return 0, errors.New("corrupted file")
	} else {
		if count, ok := paths[path]; ok {
			return count, nil
		}
		return 0, errors.New("file not found")
	}
	// b, err := os.ReadFile(path)
	// if err != nil {
	// 	return 0, err
	// }
	// wordsCount := len(strings.Fields(string(b)))
	// return wordsCount, err
}

// fan-out
func worker(ctx context.Context, jobs <-chan Job) <-chan Result {
	out := make(chan Result)
	go func() {
		defer close(out)
		for {
			select {
			case j, ok := <-jobs:
				if !ok {
					return
				}
				wordsCount, err := processFile(j.Path)
				result := Result{Path: j.Path, Words: wordsCount, Err: err}
				select {
				case out <- result:
				case <-ctx.Done():
					return
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}

// fan-in
func merge(ctx context.Context, chans ...<-chan Result) <-chan Result {
	out := make(chan Result)
	var wg sync.WaitGroup

	wg.Add(len(chans))
	for _, ch := range chans {
		go func(ch <-chan Result) {
			defer wg.Done()
			for v := range ch {
				select {
				case out <- v:
				case <-ctx.Done():
					return
				}
			}
		}(ch)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func collector(ctx context.Context, results <-chan Result) (int, error) {
	totalWordCount := 0
loop:
	for {
		select {
		case r, ok := <-results:
			if !ok {
				break loop
			}
			if r.Err != nil {
				return 0, fmt.Errorf("error: %w", r.Err)
			}
			fmt.Println(r.Path, "->", r.Words, "words")
			totalWordCount += r.Words
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}
	return totalWordCount, nil
}

func runPipeline(ctx context.Context, paths []string, workerCount int) (int, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	workerChs := make([]<-chan Result, workerCount)

	jobs := generator(ctx, paths)
	for i := range workerCount {
		workerChs[i] = worker(ctx, jobs)
	}

	out := merge(ctx, workerChs...)
	totalCount, err := collector(ctx, out)
	if err != nil {
		cancel()
	}
	return totalCount, err
}

func Pipeline() {
	paths := []string{"file_1.txt", "file_2.txt", "file_7.txt", "file_3.txt"}
	totalCount, err := runPipeline(context.Background(), paths, 4)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Total word count:", totalCount)
	}
}

func main() {
	Pipeline()
}