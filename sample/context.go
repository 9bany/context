package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	minLatency = 20
	maxLatency = 5000
	timout     = 3000
)

// very slow operation
// one cancelation signal (timeout based)
// one cancelation signal (cancelFunc based)

func main() {
	// little program that seatches flight routes
	// we are going to use a mock backend/database

	// the purpose of this is show how the context can be used
	// to propagate cancellatioin signals across
	// different go routines/ "processes" (abstractu processes)

	rootCtx := context.Background()

	ctxWithTimeout, cancel := context.WithTimeout(rootCtx, time.Duration(timout)*time.Millisecond)

	defer cancel()
	// listen for interrupt signal
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		// Now cancel
		fmt.Println("aborting due to interrupt...")
		cancel()
		os.Exit(0)
	}()

	fmt.Println("starting to search for nyc-london")
	res, err := Search(ctxWithTimeout, "nyc", "london")
	if err != nil {
		fmt.Println("got error:", err)
		return
	}
	fmt.Println("got results:", res)
}

func Search(ctx context.Context, from, to string) ([]string, error) {
	// slowSearch()
	// we need to watch for when ctx.Done() is closed
	res := make(chan []string)
	go func() {
		res <- slowSearch(from, to)
		close(res)
	}()

	// wait for 2 events
	// either of one will be result
	for {
		select {
		case dst := <-res:
			return dst, nil
		case <-ctx.Done():
			return []string{}, ctx.Err()
		}
	}

}

// is a very slow function that goes through a series
// of operations and return a clide oof strings
func slowSearch(from, to string) []string {

	rand.Seed(time.Now().Unix())
	latency := rand.Intn(maxLatency-minLatency+1) - minLatency

	fmt.Printf("started to search for %s-%s takes %dms...\n", from, to, latency)
	time.Sleep(time.Duration(latency) * time.Millisecond)

	return []string{
		from + "-" + to + "-british airways-11am",
		from + "-" + to + "-delta airlines-12am"}
}
