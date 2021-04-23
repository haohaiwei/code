package main

import (
	"io"
	"log"
	"os"
	"sync"
	"github.com/juju/ratelimit"
)

func main() {
	var wg sync.WaitGroup
	
	wg.Add(len(os.Args[1:]))
	for _, fname := range os.Args[1:] {
		go func(fname string) {
			log.Println("start", fname)
	        bucket := ratelimit.NewBucketWithRate(100 * 1024 * 1024, 100 * 1024 * 1024)
			defer wg.Done()
			f, err := os.Open(fname)
			if err != nil {
				log.Panicln(err)
			}
			stat, err := f.Stat()
			if err != nil {
				log.Panicln(err)
			}
			sr := io.NewSectionReader(f, 0, stat.Size())
			_, err = io.Copy(os.Stdout, ratelimit.Reader(sr, bucket))
			if err != nil {
				log.Panicln(err)
			}
			log.Println("done", fname)
		}(fname)
	}
	wg.Wait()
}
