package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

func computeSHA256(p string) (string, error) {
	f, err := os.Open(p)
	if err != nil {
		return "", err
	}

	defer f.Close()

	h := sha256.New()
	_, err = io.Copy(h, f)
	if err != nil {
		return "", err
	}

	hash := h.Sum(nil)

	return hex.EncodeToString(hash), nil
}

var numWorker = runtime.NumCPU()

func worker(wg *sync.WaitGroup, paths <-chan string, output chan<- string) {
	defer wg.Done()

	for p := range paths {
		hash, err := computeSHA256(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error for %v: %v\n", p, err)
			continue
		}

		output <- fmt.Sprintf("%v  %v", hash, p)
	}
}

func main() {
	var wg sync.WaitGroup

	paths := make(chan string)
	output := make(chan string)

	for i := 0; i < numWorker; i++ {
		wg.Add(1)
		go worker(&wg, paths, output)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for _, dir := range os.Args[1:] {
			err := filepath.Walk(dir, func(p string, fi os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if fi.IsDir() {
					return nil
				}

				paths <- p

				return nil
			})

			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %v\n", err)
			}
		}

		close(paths)
	}()

	go func() {
		for out := range output {
			fmt.Println(out)
		}
	}()

	wg.Wait()
}
