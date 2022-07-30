package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {

	resChan := make(chan []string)
	ctx := context.Background()

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 6*time.Second)

	defer cancel()

	request, err := http.NewRequest("GET", "http://127.0.0.1:8080", nil)

	if err != nil {
		log.Fatal(err)
	}

	request = request.WithContext(ctxWithTimeout)

	go func() {
		res, err := http.DefaultClient.Do(request)

		if err != nil {
			log.Println(err.Error())
			return
		}

		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			log.Println(err.Error())
			return
		}

		resChan <- []string{"You has data, please ! handle it"}
	}()

	for {
		select {
		case <-ctxWithTimeout.Done():
			err := ctxWithTimeout.Err()
			log.Println(err.Error())
			os.Exit(0)
		case dst := <-resChan:
			log.Println(dst)
			os.Exit(0)
		}
	}

}
