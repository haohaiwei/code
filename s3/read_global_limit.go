package main

import (
	"io"
	"log"
	"os"
	"sync"
	"strconv"
	"github.com/juju/ratelimit"
)

func main() {
	var wg sync.WaitGroup
	args := os.Args[1]
	a ,_ := strconv.ParseInt(args, 10, 64)
	b ,_ := strconv.ParseFloat(args, 64)
	
	bucket := ratelimit.NewBucketWithRate(b * 1024 * 1024, a * 1024 * 1024)
	wg.Add(len(os.Args[2:]))
	for _, fname := range os.Args[2:] {
		go func(fname string) {
			log.Println("start", fname)
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
