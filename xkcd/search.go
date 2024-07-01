package xkcd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func parseComics(num int, wg *sync.WaitGroup, comicsChan chan<- *Comics, errChan chan<- error) {
	defer wg.Done()

	resp, err := http.Get(XkcdLink + strconv.Itoa(num) + "/info.0.json")
	if err != nil {
		errChan <- err
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errChan <- fmt.Errorf("Request error: %s", resp.Status)
		return
	}

	var result Comics
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		errChan <- err
		return
	}
	comicsChan <- &result
}

// This func handles data from parseComics and saves data to .json
func SaveToFile() error {
	file, err := os.Create("comics.json")
	if err != nil {
		log.Println("Error creating file")
		return err
	}
	defer file.Close()

	var comics []*Comics
	count, errCount := 0, 0
	started := time.Now()
	fmt.Printf("Starting parsing... (%.19v)\n", started)

	// Channels
	comicsChan := make(chan *Comics)
	errChan := make(chan error)
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}

	for i := 1; i < 100; i++ {
		if i == 404 {
			continue
		}
		wg.Add(1)
		go parseComics(i, &wg, comicsChan, errChan)
	}

	// This func appends data from comicsChan to comics slice
	go func() {
		for comic := range comicsChan {
			mutex.Lock()
			comics = append(comics, comic)
			count++
			mutex.Unlock()
			//fmt.Printf("comics = %v + %v\n", comics, comic)
		}
	}()

	// Error channel
	go func() {
		for err := range errChan {
			log.Printf("Comics parsing error: %v", err)
			errCount++
			if errCount > 1 {
				break
			}
		}
	}()

	wg.Wait()
	close(comicsChan)
	close(errChan)

	// Time counting
	ended := time.Now()
	duration := ended.Sub(started)
	fmt.Printf("Parsed %d comics in %.2f sec.\n", count, duration.Seconds())

	// Writing data
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	err = enc.Encode(comics)
	if err != nil {
		return err
	}
	return nil
}
